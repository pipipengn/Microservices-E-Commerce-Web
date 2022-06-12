package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"inventory_srv/global"
	"inventory_srv/model"
	"inventory_srv/proto"
)

// SetInv 设置库存
func (s *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inventory model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inventory)
	inventory.Goods = req.GoodsId
	inventory.Stocks = req.Num
	global.DB.Save(&inventory)

	return &emptypb.Empty{}, nil
}

// InvDetail 获取库存详情
func (s *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inventory model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inventory); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "库存信息不存在")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inventory.Goods,
		Num:     inventory.Stocks,
	}, nil
}

func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// select+insert 解决幂等性: 先查询表，查看每个订单是否已经扣减过库存了===高并发漏洞===得加上分布式锁
	mutex := global.RedSync.NewMutex(fmt.Sprintf("Select_Insert_%s", uuid.NewV4()))
	if err := mutex.Lock(); err != nil {
		zap.S().Info("select+insert: redis lock 失败")
		return nil, status.Errorf(codes.Internal, "select+insert: 获取redis分布式锁异常")
	}

	var table model.StockSellDetail
	if result := global.DB.Where(&model.StockSellDetail{OrderSn: req.OrderSn}).First(&table); result.RowsAffected != 0 {
		return &emptypb.Empty{}, nil
	}

	if ok, err := mutex.Unlock(); !ok || err != nil {
		zap.S().Info("select+insert: redis unlock 失败")
		return nil, status.Errorf(codes.Internal, "select+insert: 释放redis分布式锁异常")
	}
	// =========================================================

	// 开启事务，确保每个订单一起成功或一起失败
	tx := global.DB.Begin()

	// 保存每件商品的销售情况
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var goodsDetail []model.GoodsDetail
	for _, goodInvInfo := range req.GoodsInfo {
		goodsDetail = append(goodsDetail, model.GoodsDetail{
			Goods: goodInvInfo.GoodsId,
			Num:   goodInvInfo.Num,
		})
	}
	sellDetail.Detail = goodsDetail
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建商品的销售情况失败")
	}

	// 开始一件件扣减库存
	for _, goodInvInfo := range req.GoodsInfo {
		var inventory model.Inventory
		// redis锁
		mutex := global.RedSync.NewMutex(fmt.Sprintf("goods_%d", goodInvInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			zap.S().Info("redis lock 失败")
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		// 判断是否存在数据
		if result := global.DB.Where(&model.Inventory{Goods: goodInvInfo.GoodsId}).First(&inventory); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Error(codes.NotFound, "库存信息不存在")
		}
		// 库存数小于购买数
		if inventory.Stocks < goodInvInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		// 扣减， 并发场景可能数据不一致 -- 锁，分布式锁
		inventory.Stocks -= goodInvInfo.Num
		tx.Save(&inventory)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			zap.S().Info("redis unlock 失败")
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	// 提交事务
	tx.Commit()

	return &emptypb.Empty{}, nil
}

// Reback 库存归还
func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 1.订单超时归还 2.订单创建失败 3.用户取消订单
	tx := global.DB.Begin()
	for _, goodInvInfo := range req.GoodsInfo {
		var inventory model.Inventory
		// 判断是否存在数据
		if result := global.DB.Where(&model.Inventory{Goods: goodInvInfo.GoodsId}).First(&inventory); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Error(codes.NotFound, "库存信息不存在")
		}
		// 归还， 并发场景可能数据不一致 -- 锁，分布式锁
		inventory.Stocks += goodInvInfo.Num
		tx.Save(&inventory)
	}
	// 提交事务
	tx.Commit()

	return &emptypb.Empty{}, nil
}

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Info("解析json失败")
			return consumer.ConsumeRetryLater, status.Errorf(codes.Internal, err.Error())
		}

		// ===将库存加回去，将selldetail的status改为2
		tx := global.DB.Begin()

		var sellDetail model.StockSellDetail
		if result := tx.Where(&model.StockSellDetail{
			OrderSn: orderInfo.OrderSn,
			Status:  1,
		}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}

		// 逐个归还库存
		for _, goodsDetail := range sellDetail.Detail {
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: goodsDetail.Goods}).Update("stocks", gorm.Expr("stocks+?", goodsDetail.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		// 修改销售表状态
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: sellDetail.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}

		tx.Commit()
	}

	return consumer.ConsumeSuccess, nil
}

// =================lock=================
func beiguanLock(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务，确保每个订单一起成功或一起失败
	tx := global.DB.Begin()
	for _, goodInvInfo := range req.GoodsInfo {
		var inventory model.Inventory
		// 判断是否存在数据

		//悲观锁
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInvInfo.GoodsId}).First(&inventory); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Error(codes.NotFound, "库存信息不存在")
		}

		// 库存数小于购买数
		if inventory.Stocks < goodInvInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		// 扣减， 并发场景可能数据不一致 -- 锁，分布式锁
		inventory.Stocks -= goodInvInfo.Num

		tx.Save(&inventory)
	}
	// 提交事务
	tx.Commit()

	return &emptypb.Empty{}, nil
}

func leguanLock(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务，确保每个订单一起成功或一起失败
	tx := global.DB.Begin()
	for _, goodInvInfo := range req.GoodsInfo {
		var inventory model.Inventory
		for {
			// 判断是否存在数据
			// 乐观锁
			if result := global.DB.Where(&model.Inventory{Goods: goodInvInfo.GoodsId}).First(&inventory); result.RowsAffected == 0 {
				tx.Rollback() // 回滚之前的操作
				return nil, status.Error(codes.NotFound, "库存信息不存在")
			}
			// 库存数小于购买数
			if inventory.Stocks < goodInvInfo.Num {
				tx.Rollback() // 回滚之前的操作
				return nil, status.Error(codes.ResourceExhausted, "库存不足")
			}
			// 扣减， 并发场景可能数据不一致 -- 锁，分布式锁
			// 乐观锁
			newData := map[string]interface{}{
				"stocks":  inventory.Stocks - goodInvInfo.Num,
				"version": inventory.Version + 1,
			}
			// update inventory set stock=stock-2, version=version+1 where goods=421 and version=version
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{
				Goods:   goodInvInfo.GoodsId,
				Version: inventory.Version,
			}).Updates(newData); result.RowsAffected == 0 {
				zap.S().Info("库存扣减失败")
			} else {
				break
			}
		}
	}
	// 提交事务
	tx.Commit()

	return &emptypb.Empty{}, nil
}
