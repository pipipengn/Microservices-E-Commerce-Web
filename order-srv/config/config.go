package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type RocketMQConfig struct {
	Address string `mapstructure:"address" json:"address"`
}

type JaegerConfig struct {
	Address string `mapstructure:"address" json:"address"`
}

type ServerConfig struct {
	Name             string             `mapstructure:"name" json:"name"`
	MysqlInfo        MysqlConfig        `mapstructure:"mysql" json:"mysql"`
	ConsulInfo       ConsulConfig       `mapstructure:"consul" json:"consul"`
	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	Debug            bool               `mapstructure:"debug"`
	Tags             []string           `mapstructure:"tags" json:"tags"`
	Host             string             `mapstructure:"host" json:"host"`
	RedisInfo        RedisConfig        `mapstructure:"redis" json:"redis"`
	RocketMQInfo     RocketMQConfig     `mapstructure:"rocketmq" json:"rocketmq"`
	JaegerInfo       JaegerConfig       `mapstructure:"jaeger" json:"jaeger"`
}
