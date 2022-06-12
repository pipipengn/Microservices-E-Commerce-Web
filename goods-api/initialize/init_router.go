package initialize

import (
	"github.com/gin-gonic/gin"
	"goods-web/api/goods"
	"goods-web/middlewares"
	"goods-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())

	Router.GET("/s3Url", goods.GetPresignedS3Url)
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	ApiGroup := Router.Group("/g/v1")
	ApiGroup.Use(middlewares.Trace())

	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	return Router
}
