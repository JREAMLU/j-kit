package http

import (
	"log"
	"os"
	"time"

	middleware "github.com/JREAMLU/j-kit/go-micro/gin-middleware"
	microGobreaker "github.com/JREAMLU/j-kit/go-micro/plugins/wrapper/breaker/gobreaker"
	jopentracing "github.com/JREAMLU/j-kit/go-micro/trace/opentracing"
	"github.com/JREAMLU/j-kit/go-micro/util"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	micro "github.com/micro/go-micro"
	microClient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	client "github.com/micro/go-plugins/client/grpc"
	register "github.com/micro/go-plugins/registry/consul"
	server "github.com/micro/go-plugins/server/grpc"
	transport "github.com/micro/go-plugins/transport/grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
)

// NewHTTPService new http service
func NewHTTPService(config *util.Config) (micro.Service, *gin.Engine, opentracing.Tracer) {
	t, err := util.NewTrace(
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
			microClient.Wrap(microGobreaker.NewClientWrapper(circuitBreakers(config))),
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
		micro.WrapClient(jopentracing.NewClientWrapper(t)),
		micro.WrapHandler(jopentracing.NewHandlerWrapper(t)),
	)

	g := gin.New()
	g.Use(
		gin.Recovery(),
		HandlerHTTPRequestGin(t, config.Service.Name),
		middleware.HeaderTraceRespone(),
	)

	return service, g, t
}

func circuitBreakers(config *util.Config) map[string]*gobreaker.CircuitBreaker {
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
