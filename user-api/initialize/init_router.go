package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-web/middlewares"
	"user-web/router"
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

	ApiGroup := Router.Group("/u/v1")

	router.InitBaseRouter(ApiGroup)
	router.InitUserRouter(ApiGroup)

	return Router
}
