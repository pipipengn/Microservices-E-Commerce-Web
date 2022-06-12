package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"order_srv/global"
	"order_srv/handler"
	"order_srv/initialize"
	"order_srv/proto"
	"order_srv/utils"
	"order_srv/utils/otgrpc"
	"order_srv/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()

	zap.S().Info(global.ServerConfig)

	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 50054, "port")

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info(*IP)
	zap.S().Info(*Port)

	// 初始化jaeger
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: global.ServerConfig.JaegerInfo.Address,
		},
		ServiceName: "shop",
	}

	// new tracer
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(tracer)
	if err != nil {
		panic(err)
	}

	// 初始化grpc
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterOrderServer(server, &handler.OrderServer{})
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

	// 监听库存归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{global.ServerConfig.RocketMQInfo.Address}),
		consumer.WithGroupName("shop-order"),
	)

	// 订阅mq的消息
	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		zap.S().Error("读取mq消息失败", err.Error())
	}
	_ = c.Start()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = c.Shutdown()
	_ = closer.Close()
	if err = registryClient.DeRegister(serviceId); err != nil {
		zap.S().Info("服务注销失败", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
