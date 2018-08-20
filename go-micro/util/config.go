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
		// when StateHalfOpen, allow how many requests try in
		MaxRequests uint32
		// when StateClosed, after every interval time, clean Counts.Requests (failure requests)
		Interval int
		// when StateOpen, after every timeout, change to StateHalfOpen
		Timeout int
		// in interval counts requests
		CountsRequests uint32
		// failure ratio
		FailureRatio float64
	}

	RateLimit struct {
		// gvie bucket rate time (every time give bucket nums)
		ClientRate float64
		// total bucket
		ClientCapacity int64
		ClientWait     bool
		ServerRate     float64
		ServerCapacity int64
		ServerWait     bool
	}

	Kafka struct {
		Broker        []string
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
	_broker                 = "broker"
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

		config.Kafka.Broker, _, err = client.GetKafkas(_broker)
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
		config.CircuitBreaker.MaxRequests = 100
	}

	if config.CircuitBreaker.FailureRatio == 0 {
		config.CircuitBreaker.FailureRatio = 0.6
	}

	if config.CircuitBreaker.Interval == 0 {
		config.CircuitBreaker.Interval = 30
	}

	if config.CircuitBreaker.Timeout == 0 {
		config.CircuitBreaker.Timeout = 90
	}

	if config.CircuitBreaker.CountsRequests == 0 {
		config.CircuitBreaker.CountsRequests = 1000
	}

	if config.RateLimit.ClientRate == 0 {
		config.RateLimit.ClientRate = 2000
	}

	if config.RateLimit.ClientCapacity == 0 {
		config.RateLimit.ClientCapacity = 10000
	}

	if config.RateLimit.ServerRate == 0 {
		config.RateLimit.ServerRate = 2000
	}

	if config.RateLimit.ServerCapacity == 0 {
		config.RateLimit.ServerCapacity = 10000
	}

	log.Printf("Zipkin Broker: %v", config.Kafka.ZipkinBroker)
	log.Printf("Zipkin Topic: %v", config.Kafka.ZipkinTopic)

	log.Printf("Broker Broker: %v", config.Kafka.Broker)

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

	log.Printf("RateLimit ClientRate: %v", config.RateLimit.ClientRate)
	log.Printf("RateLimit ClientCapacity: %v", config.RateLimit.ClientCapacity)
	log.Printf("RateLimit ClientWait: %v", config.RateLimit.ClientWait)
	log.Printf("RateLimit ServerRate: %v", config.RateLimit.ServerRate)
	log.Printf("RateLimit ServerCapacity: %v", config.RateLimit.ServerCapacity)
	log.Printf("RateLimit ServerWait: %v", config.RateLimit.ServerWait)

	return nil
}

func getServiceKey(name, version string) string {
	return ext.StringSplice(serviceGo, name, "/", version)
}
