package handler

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order_srv/global"
	"order_srv/model"
	"order_srv/proto"
	"time"
)

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	parentSpan := opentracing.SpanFromContext(o.Ctx)

	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	// 从购物车中获取到选中的商品
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	var buyingCart []model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&buyingCart); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	// 拿到商品id和每件商品的数量
	var goodsIds []int32
	goodsNumsMap := make(map[int32]int32)
	for _, buyingGoods := range buyingCart {
		goodsIds = append(goodsIds, buyingGoods.Goods)
		goodsNumsMap[buyingGoods.Goods] = buyingGoods.Nums
	}

	// 商品的价格自己查询 - 访问商品服务 (跨微服务)
	goodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goodsList, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	goodsSpan.Finish()

	// 计算总价格 + 订单的商品信息表
	var orderPrice float32
	var orderGoods []*model.OrderGoods
	for _, goods := range goodsList.Data {
		orderPrice += goods.ShopPrice * float32(goodsNumsMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       goodsNumsMap[goods.Id],
		})
	}

	// 库存的扣减 - 访问库存服务 (跨微服务)
	var goodsInfo []*proto.GoodsInvInfo
	for _, buyingGoods := range buyingCart {
		goodsInfo = append(goodsInfo, &proto.GoodsInvInfo{
			GoodsId: buyingGoods.Goods,
			Num:     buyingGoods.Nums,
		})
	}
	InventorySpan := opentracing.GlobalTracer().StartSpan("query_inventorty", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInfo,
		OrderSn:   orderInfo.OrderSn,
	}); err != nil {
		// 如果因为网络问题，其实扣减成功：判断sell返回的error状态码
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	InventorySpan.Finish()

	// 开启本地事务
	tx := global.DB.Begin()
	// 订单的基本信息表
	orderInfo.OrderMount = orderPrice
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	// 订单的商品信息表中添加订单id, 并保存到数据库
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_orderGoods", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单商品信息失败"
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	// 删除购物车记录
	deleteShopcartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	deleteShopcartSpan.Finish()

	// 延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.ServerConfig.RocketMQInfo.Address}))
	if err != nil {
		panic(err)
	}

	if err = p.Start(); err != nil {
		panic(err)
	}

	message := primitive.NewMessage("order_timeout", msg.Body)
	if _, err = p.SendSync(context.Background(), message.WithDelayTimeLevel(3)); err != nil {
		zap.S().Error("发送延时消息失败")
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.RollbackMessageState
	}

	tx.Commit()

	o.OrderAmount = orderPrice
	o.ID = orderInfo.ID
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	// 检查之前逻辑是否完成
	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState // 并不能说明库存已经扣减，在另一端要做好幂等性的保证
	}
	return primitive.RollbackMessageState
}

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/

	orderListener := OrderListener{Ctx: ctx}
	p, err := rocketmq.NewTransactionProducer(&orderListener, producer.WithNameServer([]string{global.ServerConfig.RocketMQInfo.Address}))
	if err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}

	if err = p.Start(); err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}

	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	bytes, _ := json.Marshal(order)

	_, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", bytes))
	if err != nil {
		return nil, status.Error(codes.Internal, "消息发送失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderAmount,
	}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		var orderInfo OrderInfo
		if err := json.Unmarshal(msgs[i].Body, &orderInfo); err != nil {
			zap.S().Info("解析json失败")
			return consumer.ConsumeRetryLater, status.Errorf(codes.Internal, err.Error())
		}

		zap.S().Info("订单超时", time.Now())
		// 查询订单状态，如果以支付，什么都不做，未支付，归还库存
		var order model.OrderInfo
		if result := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}

		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()

			// 修改订单状态为取消
			order.Status = "TRADE_CANCLE"
			if result := tx.Save(&order); result.RowsAffected == 0 {
				return consumer.ConsumeRetryLater, nil
			}

			// 归还库存，发个消息到order_reback中去
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.ServerConfig.RocketMQInfo.Address}))
			if err != nil {
				tx.Rollback()
				zap.S().Error("连接mq失败")
				return consumer.ConsumeRetryLater, nil
			}

			if err = p.Start(); err != nil {
				tx.Rollback()
				zap.S().Error("启动mq失败")
				return consumer.ConsumeRetryLater, nil
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				zap.S().Error("发送消息失败")
				return consumer.ConsumeRetryLater, nil
			}
			tx.Commit()
		}
	}

	return consumer.ConsumeSuccess, nil
}
