package util

import (
	"log"
	"os"
	"time"

	microGobreaker "github.com/JREAMLU/j-kit/go-micro/plugins/wrapper/breaker/gobreaker"
	microRatelimit "github.com/JREAMLU/j-kit/go-micro/plugins/wrapper/ratelimiter/ratelimit"
	"github.com/JREAMLU/j-kit/go-micro/trace/opentracing"
	"github.com/juju/ratelimit"

	"github.com/hashicorp/consul/api"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	microClient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	microServer "github.com/micro/go-micro/server"
	brokerKafka "github.com/micro/go-plugins/broker/kafka"
	client "github.com/micro/go-plugins/client/grpc"
	register "github.com/micro/go-plugins/registry/consul"
	server "github.com/micro/go-plugins/server/grpc"
	transport "github.com/micro/go-plugins/transport/grpc"
	"github.com/sony/gobreaker"
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

	serverBucket := ratelimit.NewBucketWithRate(config.ServerRateLimit.Rate, config.ServerRateLimit.Capacity)

	service := micro.NewService(
		micro.Client(client.NewClient(
			microClient.Wrap(microGobreaker.NewClientWrapper(circuitBreakers(config))),
			microClient.Wrap(microRatelimit.NewClientWrapper(clientBuckets(config))),
		)),
		micro.Server(server.NewServer(
			microServer.WrapHandler(microRatelimit.NewHandlerWrapper(serverBucket, config.ServerRateLimit.Wait)),
		)),
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
		micro.Broker(brokerKafka.NewBroker(
			broker.Option(func(opt *broker.Options) {
				opt.Addrs = config.Kafka.Broker
			}),
		)),
	)

	service.Init(
		micro.RegisterTTL(time.Duration(config.Service.RegisterTTL)*time.Second),
		micro.RegisterInterval(time.Duration(config.Service.RegisterInterval)*time.Second),
	)

	return service
}

func circuitBreakers(config *Config) map[string]*gobreaker.CircuitBreaker {
	cbs := make(map[string]*gobreaker.CircuitBreaker, len(config.CircuitBreakers))

	for circuitName, circuitBreaker := range config.CircuitBreakers {
		cbs[circuitName] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        circuitName,
			MaxRequests: circuitBreaker.MaxRequests,
			Interval:    time.Duration(circuitBreaker.Interval) * time.Second,
			Timeout:     time.Duration(circuitBreaker.Timeout) * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= circuitBreaker.CountsRequests && failureRatio >= circuitBreaker.FailureRatio
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			},
		})

	}

	return cbs
}

func clientBuckets(config *Config) (map[string]*ratelimit.Bucket, bool) {
	bs := make(map[string]*ratelimit.Bucket, len(config.ClientRateLimits))
	var wait bool

	for rateName, rateLimit := range config.ClientRateLimits {
		bs[rateName] = ratelimit.NewBucketWithRate(rateLimit.Rate, rateLimit.Capacity)
		wait = rateLimit.Wait
	}

	return bs, wait
}
