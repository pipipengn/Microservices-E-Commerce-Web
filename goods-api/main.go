package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"goods-web/global"
	"goods-web/initialize"
	"goods-web/utils"
	"goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()
	initialize.InitSentinel()
	r := initialize.Routers()

	if !global.ServerConfig.Debug {
		global.ServerConfig.Port, _ = utils.GetFreePort()
	}

	// 服务注册
	registryClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serverConfig := global.ServerConfig
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registryClient.Register(serverConfig.Host, serverConfig.Port, serverConfig.Name, serverConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}

	zap.S().Infof("服务器启动端口: %d", global.ServerConfig.Port)
	go func() {
		// 启动服务
		if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("server start fail", err.Error())
		}
	}()

	// 退出时注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = registryClient.DeRegister(serviceId); err != nil {
		zap.S().Info("服务注销失败", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
