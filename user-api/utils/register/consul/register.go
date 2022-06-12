package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"user-web/global"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registery{
		Host: host,
		Port: port,
	}
}

type Registery struct {
	Host string
	Port int
}

func (r *Registery) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.Host, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Address = address
	registration.Port = port
	registration.Tags = tags
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registery) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		return err
	}
	return nil
}
