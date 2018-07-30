package util

import (
	"errors"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/core/consul"
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
		BigdataBroker []string
	}

	Zookeeper struct {
		BigdataAddrs  []string
		BigdataZkroot string
	}
}

const (
	serviceGo = "service/go/"
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
		// @TODO kafka zookeeper
	}

	return nil
}

func getServiceKey(name, version string) string {
	return ext.StringSplice(serviceGo, name, "/", version)
}
