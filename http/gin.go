package http

import (
	"log"
	"os"

	jopentracing "github.com/JREAMLU/j-kit/go-micro/trace/opentracing"
	"github.com/JREAMLU/j-kit/go-micro/util"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	client "github.com/micro/go-plugins/client/grpc"
	register "github.com/micro/go-plugins/registry/consul"
	server "github.com/micro/go-plugins/server/grpc"
	transport "github.com/micro/go-plugins/transport/grpc"
	opentracing "github.com/opentracing/opentracing-go"
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
		micro.Client(client.NewClient()),
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
	)

	return service, g, t
}
