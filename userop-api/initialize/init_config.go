package initialize

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"userop-web/global"
)

var (
	host      = "172.31.53.72"
	port      = 8848
	namespace = "53d73c25-af57-4007-917f-ebff62942b92"
	dataid    = "userop-web"
	group     = "prod"
)

func InitConfig() {
	// 从nacos中读其他配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: host,
			Port:   uint64(port),
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// Create config client for dynamic configuration
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataid,
		Group:  group,
	})
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos失败 %s", err)
	}
}
