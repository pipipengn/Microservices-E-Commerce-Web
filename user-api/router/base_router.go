package router

import (
	"github.com/gin-gonic/gin"
	"user-web/api"
)

func InitBaseRouter(router *gin.RouterGroup) {

	BaseRouter := router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("sms", api.SendSms)
	}
}
