package pay

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"order-web/global"
	"order-web/proto"
)

func Notify(c *gin.Context) {

	orderSn := c.PostForm("order_sn")
	status := c.PostForm("status")

	if _, err := global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: orderSn,
		Status:  status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.String(http.StatusOK, "success")
}
