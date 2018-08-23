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

	CircuitBreakers map[string]struct {
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

	ClientRateLimits map[string]struct {
		// gvie bucket rate time (every time give bucket nums)
		Rate float64
		// total bucket
		Capacity int64
		Wait     bool
	}

	ServerRateLimit struct {
		Rate     float64
		Capacity int64
		Wait     bool
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
	for circuitName, circuitBreaker := range config.CircuitBreakers {
		if circuitBreaker.MaxRequests == 0 {
			circuitBreaker.MaxRequests = 100
		}

		if circuitBreaker.FailureRatio == 0 {
			circuitBreaker.FailureRatio = 0.6
		}

		if circuitBreaker.Interval == 0 {
			circuitBreaker.Interval = 30
		}

		if circuitBreaker.Timeout == 0 {
			circuitBreaker.Timeout = 90
		}

		if circuitBreaker.CountsRequests == 0 {
			circuitBreaker.CountsRequests = 1000
		}

		log.Printf("%v CircuitBreaker MaxRequests: %v, FailureRatio: %v, Interval: %v, Timeout: %v, CountsRequests: %v",
			circuitName, circuitBreaker.MaxRequests, circuitBreaker.FailureRatio, circuitBreaker.Interval, circuitBreaker.Timeout, circuitBreaker.CountsRequests)
	}

	for rateName, rateLimit := range config.ClientRateLimits {
		if rateLimit.Rate == 0 {
			rateLimit.Rate = 2000
		}

		if rateLimit.Capacity == 0 {
			rateLimit.Capacity = 10000
		}

		log.Printf("%v RateLimit ClientRate: %v, ClientCapacity: %v, ClientWait: %v", rateName, rateLimit.Rate, rateLimit.Capacity, rateLimit.Wait)
	}

	if config.ServerRateLimit.Rate == 0 {
		config.ServerRateLimit.Rate = 2000
	}

	if config.ServerRateLimit.Capacity == 0 {
		config.ServerRateLimit.Capacity = 10000
	}

	log.Printf("RateLimit ServerRate: %v, ServerCapacity: %v, ServerWait: %v", config.ServerRateLimit.Rate, config.ServerRateLimit.Capacity, config.ServerRateLimit.Wait)
	log.Printf("Zipkin Broker: %v, Topic: %v", config.Kafka.ZipkinBroker, config.Kafka.ZipkinTopic)
	log.Printf("Broker Broker: %v", config.Kafka.Broker)
	log.Printf("Service RegisterInterval: %v, RegisterTTL: %v", config.Service.RegisterInterval, config.Service.RegisterTTL)

	if config.Web.Port != 0 {
		log.Printf("Web Host: %v Port: %v", config.Web.Host, config.Web.Port)
	}

	return nil
}

func getServiceKey(name, version string) string {
	return ext.StringSplice(serviceGo, name, "/", version)
}
