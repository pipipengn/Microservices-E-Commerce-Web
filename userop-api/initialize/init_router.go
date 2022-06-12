package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userop-web/middlewares"
	"userop-web/router"
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

	ApiGroup := Router.Group("/up/v1")
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)

	return Router
}
