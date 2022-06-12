package global

import (
	"userop-web/config"
	"userop-web/proto"
)

var (
	ServerConfig   *config.ServerSrvConfig = &config.ServerSrvConfig{}
	NacosConfig    *config.NacosConfig     = &config.NacosConfig{}
	MessageClient  proto.MessageClient
	AddressClient  proto.AddressClient
	UserFavClient  proto.UserFavClient
	GoodsSrvClient proto.GoodsClient
)
