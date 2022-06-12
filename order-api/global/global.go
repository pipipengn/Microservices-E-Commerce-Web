package global

import (
	"order-web/config"
	"order-web/proto"
)

var (
	ServerConfig       *config.ServerSrvConfig = &config.ServerSrvConfig{}
	NacosConfig        *config.NacosConfig     = &config.NacosConfig{}
	OrderSrvClient     proto.OrderClient
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
