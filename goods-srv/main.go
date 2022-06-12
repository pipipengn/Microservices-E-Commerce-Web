package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"goods_srv/global"
	"goods_srv/handler"
	"goods_srv/initialize"
	"goods_srv/proto"
	"goods_srv/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info(global.ServerConfig)

	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 50052, "port")

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info(*IP)
	zap.S().Info(*Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("fail to listen:" + err.Error())
	}

	// 健康检查
	grpchealth.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceId
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		// 启动服务
		err = server.Serve(listen)
		if err != nil {
			panic("fail to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
