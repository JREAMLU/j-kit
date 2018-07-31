package util

import (
	"errors"
	"log"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/consul"
	"github.com/JREAMLU/j-kit/ext"
)

// Config service config
type Config struct {
	Service struct {
		Name             string
		Version          string
		RegisterTTL      int
		RegisterInterval int
	}

	Web struct {
		Host string
		Port int
	}

	Kafka struct {
		ZipkinBroker  []string
		ZipkinTopic   string
		BigdataBroker []string
	}

	Zookeeper struct {
		BigdataAddrs  []string
		BigdataZkroot string
	}

	Consul struct {
		RegistryAddrs []string
	}
}

const (
	serviceGo   = "service/go/"
	_zipkin     = "zipkin"
	_bigdata    = "bigdata"
	_consul     = "registry"
	zipkinTopic = "zipkin"
)

// LoadConfig load service config
func LoadConfig(consulAddr string, name, version string) (*Config, error) {
	sc := &Config{}
	return sc, loadConfig(consulAddr, getServiceKey(name, version), sc)
}

// LoadCustomConfig load service config by custom
func LoadCustomConfig(consulAddr string, name, version string, sc interface{}) error {
	return loadConfig(consulAddr, getServiceKey(name, version), sc)
}

func loadConfig(consulAddr string, key string, sc interface{}) error {
	client, err := consul.NewClient(consul.SetAddress(consulAddr))
	if err != nil {
		return err
	}

	buf, err := client.Get(key)
	if err != nil {
		return err
	}

	log.Printf("Load Config: %v \n%v\n%v\n%v\n\n", key, consul.SeparatorStart, buf, consul.SeparatorEnd)

	_, err = toml.Decode(buf, sc)
	if err != nil {
		return err
	}

	config, ok := sc.(*Config)
	if !ok {
		s := reflect.ValueOf(sc).Elem()
		v := s.FieldByName("Config")
		if v.IsValid() {
			config, ok = (v.Interface()).(*Config)
			if !ok {
				return errors.New("type interface not ok")
			}
		}
	}

	if config != nil {
		config.Kafka.ZipkinBroker, err = client.GetKafkas(_zipkin)
		if err != nil {
			return err
		}

		config.Zookeeper.BigdataAddrs, config.Zookeeper.BigdataZkroot, err = client.GetZookeepers(_bigdata)
		if err != nil {
			return err
		}

		config.Kafka.ZipkinTopic = zipkinTopic

		config.Consul.RegistryAddrs, err = client.GetConsulAddrs(_consul)
		if err != nil {
			return err
		}
	}

	if config.Service.RegisterInterval == 0 {
		config.Service.RegisterInterval = 1
	}

	if config.Service.RegisterTTL == 0 {
		config.Service.RegisterTTL = 1
	}

	return nil
}

func getServiceKey(name, version string) string {
	return ext.StringSplice(serviceGo, name, "/", version)
}
