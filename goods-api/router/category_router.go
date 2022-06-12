package router

import (
	"github.com/gin-gonic/gin"
	"goods-web/api/category"
	"goods-web/middlewares"
)

//func InitCategoryRouter(router *gin.RouterGroup) {
//
//	categoryRouter := router.Group("category")
//	zap.S().Info("配置category url")
//	{
//		categoryRouter.GET("/list", category.List)
//		categoryRouter.GET("/detail/:id", category.Detail)
//		categoryRouter.POST("/create", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.New)
//		categoryRouter.POST("/delete/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Delete)
//		categoryRouter.POST("/update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Update)
//	}
//}

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys").Use(middlewares.Trace())
	{
		CategoryRouter.GET("", category.List)          // 商品类别列表页
		CategoryRouter.DELETE("/:id", category.Delete) // 删除分类
		CategoryRouter.GET("/:id", category.Detail)    // 获取分类详情
		CategoryRouter.POST("", category.New)          //新建分类
		CategoryRouter.PUT("/:id", category.Update)    //修改分类信息
	}
}
