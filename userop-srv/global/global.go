package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"userop_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	RedSync      *redsync.Redsync     = &redsync.Redsync{}
)
