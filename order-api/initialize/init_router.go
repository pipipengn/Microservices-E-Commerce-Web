package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-web/api/pay"
	"order-web/middlewares"
	"order-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	Router.POST("/pay/notify", pay.Notify)

	ApiGroup := Router.Group("/o/v1")
	ApiGroup.Use(middlewares.Trace())

	router.InitOrderRouter(ApiGroup)
	router.InitShopCart(ApiGroup)

	return Router
}
