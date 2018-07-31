package util

import (
	"time"

	"github.com/JREAMLU/j-kit/go-micro/trace/opentracing"

	micro "github.com/micro/go-micro"
	clientGrpc "github.com/micro/go-plugins/client/grpc"
	registerConsul "github.com/micro/go-plugins/registry/consul"
	serverGrpc "github.com/micro/go-plugins/server/grpc"
	transportGrpc "github.com/micro/go-plugins/transport/grpc"
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
		micro.Client(clientGrpc.NewClient()),
		micro.Server(serverGrpc.NewServer()),
		micro.Registry(registerConsul.NewRegistry()),
		micro.Transport(transportGrpc.NewTransport()),
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
