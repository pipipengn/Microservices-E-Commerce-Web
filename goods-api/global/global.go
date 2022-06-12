package global

import (
	"goods-web/config"
	"goods-web/proto"
)

var (
	ServerConfig   *config.ServerSrvConfig = &config.ServerSrvConfig{}
	NacosConfig    *config.NacosConfig     = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
)
