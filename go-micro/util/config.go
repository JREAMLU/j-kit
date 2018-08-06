package util

import (
	"errors"
	"fmt"
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
		URL  string
	}

	CircuitBreaker struct {
		MaxRequests    uint32
		Interval       int
		Timeout        int
		CountsRequests uint32
		FailureRatio   float64
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
	serviceGo               = "service/go/"
	_zipkin                 = "zipkin"
	_bigdata                = "bigdata"
	_consul                 = "registry"
	zipkinTopic             = "zipkin"
	defaultRegisterInterval = 20
	defaultRegisterTTL      = 30
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

	// kafka zookeeper
	if config != nil {
		config.Kafka.ZipkinBroker, config.Kafka.ZipkinTopic, err = client.GetKafkas(_zipkin)
		if err != nil {
			return err
		}

		config.Zookeeper.BigdataAddrs, config.Zookeeper.BigdataZkroot, err = client.GetZookeepers(_bigdata)
		if err != nil {
			return err
		}

		if config.Kafka.ZipkinTopic == "" {
			config.Kafka.ZipkinTopic = zipkinTopic
		}

		config.Consul.RegistryAddrs, err = client.GetConsulAddrs(_consul)
		if err != nil {
			return err
		}
	}

	// micro
	if config.Service.RegisterInterval == 0 {
		config.Service.RegisterInterval = defaultRegisterInterval
	}

	if config.Service.RegisterTTL == 0 {
		config.Service.RegisterTTL = defaultRegisterTTL
	}

	// http server
	if config.Web.URL == "" {
		config.Web.URL = fmt.Sprintf("%v:%v", config.Web.Host, config.Web.Port)
	}

	// circuit Breaker
	if config.CircuitBreaker.MaxRequests == 0 {
		// when StateHalfOpen, allow how many requests try in
		config.CircuitBreaker.MaxRequests = 100
	}

	if config.CircuitBreaker.FailureRatio == 0 {
		// failure ratio
		config.CircuitBreaker.FailureRatio = 0.6
	}

	if config.CircuitBreaker.Interval == 0 {
		// when StateClosed, after every interval time, clean Counts.Requests (failure requests)
		config.CircuitBreaker.Interval = 30
	}

	if config.CircuitBreaker.Timeout == 0 {
		// when StateOpen, after every timeout, change to StateHalfOpen
		config.CircuitBreaker.Timeout = 90
	}

	if config.CircuitBreaker.CountsRequests == 0 {
		// failure requests
		config.CircuitBreaker.CountsRequests = 1000
	}

	log.Printf("Zipkin Broker: %v", config.Kafka.ZipkinBroker)
	log.Printf("Zipkin Topic: %v", config.Kafka.ZipkinTopic)
	log.Printf("Service RegisterInterval: %v", config.Service.RegisterInterval)
	log.Printf("Service RegisterTTL: %v", config.Service.RegisterTTL)
	if config.Web.Port != 0 {
		log.Printf("Web Host: %v Port: %v", config.Web.Host, config.Web.Port)
	}
	log.Printf("CircuitBreaker MaxRequests: %v", config.CircuitBreaker.MaxRequests)
	log.Printf("CircuitBreaker FailureRatio: %v", config.CircuitBreaker.FailureRatio)
	log.Printf("CircuitBreaker Interval: %v", config.CircuitBreaker.Interval)
	log.Printf("CircuitBreaker Timeout: %v", config.CircuitBreaker.Timeout)
	log.Printf("CircuitBreaker CountsRequests: %v", config.CircuitBreaker.CountsRequests)

	return nil
}

func getServiceKey(name, version string) string {
	return ext.StringSplice(serviceGo, name, "/", version)
}
