package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"user-web/api"
	"user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	router.POST("pwdlogin", api.PasswordLogin)
	router.POST("register", api.Register)

	UserRouter := router.Group("user").Use(middlewares.JWTAuth(), middlewares.IsAdminAuth())
	zap.S().Infof("配置用户URL")
	{
		UserRouter.GET("", api.GetUserList)
	}
}
