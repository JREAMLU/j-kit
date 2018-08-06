package util

import (
	"log"
	"os"
	"time"

	"github.com/JREAMLU/j-kit/go-micro/trace/opentracing"

	"github.com/hashicorp/consul/api"
	micro "github.com/micro/go-micro"
	microClient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	client "github.com/micro/go-plugins/client/grpc"
	register "github.com/micro/go-plugins/registry/consul"
	server "github.com/micro/go-plugins/server/grpc"
	transport "github.com/micro/go-plugins/transport/grpc"
	microGobreaker "github.com/micro/go-plugins/wrapper/breaker/gobreaker"
	"github.com/sony/gobreaker"
	// brokerKafka "github.com/micro/go-plugins/broker/kafka"
)

// NewMicroService new micro service
func NewMicroService(config *Config) micro.Service {
	t, err := NewTrace(
		config.Service.Name,
		config.Service.Version,
		config.Kafka.ZipkinBroker,
		config.Kafka.ZipkinTopic,
	)
	if err != nil {
		panic(err)
	}

	service := micro.NewService(
		micro.Client(client.NewClient(
			microClient.Wrap(microGobreaker.NewClientWrapper(
				gobreaker.NewCircuitBreaker(gobreaker.Settings{
					Name:        config.Service.Name,
					MaxRequests: config.CircuitBreaker.MaxRequests,
					Interval:    time.Duration(config.CircuitBreaker.Interval) * time.Second,
					Timeout:     time.Duration(config.CircuitBreaker.Timeout) * time.Second,
					ReadyToTrip: func(counts gobreaker.Counts) bool {
						failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
						return counts.Requests >= config.CircuitBreaker.CountsRequests && failureRatio >= config.CircuitBreaker.FailureRatio
					},
					OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
					},
				}),
			)),
		)),
		micro.Server(server.NewServer()),
		micro.Registry(register.NewRegistry(
			registry.Option(func(opts *registry.Options) {
				if len(config.Consul.RegistryAddrs) > 0 {
					opts.Addrs = config.Consul.RegistryAddrs
					log.Printf("Registry Consul Addrs: %v\n", config.Consul.RegistryAddrs)
				} else {
					log.Printf("Registry Consul Addr: %v\n", os.Getenv(api.HTTPAddrEnvName))
				}
			}),
		)),
		micro.Transport(transport.NewTransport()),
		micro.Name(config.Service.Name),
		micro.Version(config.Service.Version),
		micro.WrapClient(opentracing.NewClientWrapper(t)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		// micro.Broker(brokerKafka.NewBroker(
		// 	broker.Option(func(opt *broker.Options) {
		// 		opt.Addrs = []string{"10.200.119.128:9092"}
		// 	}),
		// )),
	)

	service.Init(
		micro.RegisterTTL(time.Duration(config.Service.RegisterTTL)*time.Second),
		micro.RegisterInterval(time.Duration(config.Service.RegisterInterval)*time.Second),
	)

	return service
}
