package shopcart

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"order-web/api"
	"order-web/forms"
	"order-web/global"
	"order-web/proto"
	"strconv"
)

type ListReturnInfo struct {
	Id      int32
	Nums    int32
	checked bool
}

func List(c *gin.Context) {
	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	ids := make([]int32, 0)
	returnInfoMap := make(map[int32]ListReturnInfo, 0)

	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
		returnInfoMap[item.GoodsId] = ListReturnInfo{
			Id:      item.Id,
			Nums:    item.Nums,
			checked: item.Checked,
		}
	}
	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// CartItemList的返回值只包含goodsId   请求商品服务获取商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	/*
		{
			"total":12,
			"data":[
				{
					"id":1,
					"goods_id":421,
					"goods_name":421,
					"goods_price":421,
					"goods_image":421,
					"nums":421,
					"checked":421,
				}
			]
		}
	*/
	goodsList := make([]interface{}, 0)
	for _, good := range goodsRsp.Data {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = returnInfoMap[good.Id].Id
		tmpMap["goods_id"] = good.Id
		tmpMap["goods_name"] = good.Name
		tmpMap["goods_image"] = good.GoodsFrontImage
		tmpMap["goods_price"] = good.ShopPrice
		tmpMap["nums"] = returnInfoMap[good.Id].Nums
		tmpMap["checked"] = returnInfoMap[good.Id].checked

		goodsList = append(goodsList, tmpMap)
	}

	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//为了严谨性，添加商品到购物车之前，检查商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("[New] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	//如果现在添加到购物车的数量和库存的数量不一致
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("[New] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	if invRsp.Num < itemForm.Nums {
		c.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("userId")
	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	userId, _ := c.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}
