package global

import (
	"user-web/config"
	"user-web/proto"
)

var (
	ServerConfig  *config.ServerSrvConfig = &config.ServerSrvConfig{}
	NacosConfig   *config.NacosConfig     = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
