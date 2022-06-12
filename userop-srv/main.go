package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"userop_srv/global"
	"userop_srv/handler"
	"userop_srv/initialize"
	"userop_srv/proto"
	"userop_srv/utils"
	"userop_srv/utils/register/consul"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedSync()

	zap.S().Info(global.ServerConfig)

	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 50055, "port")

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info(*IP)
	zap.S().Info(*Port)

	// 初始化grpc
	server := grpc.NewServer()
	proto.RegisterMessageServer(server, &handler.UserOpServer{})
	proto.RegisterAddressServer(server, &handler.UserOpServer{})
	proto.RegisterUserFavServer(server, &handler.UserOpServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("fail to listen:" + err.Error())
	}

	// 健康检查
	grpchealth.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	registryClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serverConfig := global.ServerConfig
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = registryClient.Register(serverConfig.Host, *Port, serverConfig.Name, serverConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}

	zap.S().Infof("服务器启动端口: %d", *Port)

	// 启动服务
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("fail to start grpc:" + err.Error())
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = registryClient.DeRegister(serviceId); err != nil {
		zap.S().Info("服务注销失败", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
