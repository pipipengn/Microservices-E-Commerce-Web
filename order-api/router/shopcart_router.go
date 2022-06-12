package router

import (
	"github.com/gin-gonic/gin"
	"order-web/api/shopcart"
	"order-web/middlewares"
)

func InitShopCart(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", shopcart.List)          //购物车列表
		GoodsRouter.DELETE("/:id", shopcart.Delete) //删除条目
		GoodsRouter.POST("", shopcart.New)          //添加商品到购物车
		GoodsRouter.PATCH("/:id", shopcart.Update)  //修改条目
	}
}
