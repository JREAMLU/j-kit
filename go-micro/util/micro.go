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
func NewMicroService() micro.Service {
	t, err := NewTrace("go.micro.srv.s1", "v1", []string{"10.200.119.128:9092", "10.200.119.129:9092", "10.200.119.130:9092"}, "web_log_get")
	if err != nil {
		panic(err)
	}

	service := micro.NewService(
		micro.Client(clientGrpc.NewClient()),
		micro.Server(serverGrpc.NewServer()),
		micro.Registry(registerConsul.NewRegistry()),
		micro.Transport(transportGrpc.NewTransport()),
		micro.Name("go.micro.srv.s1"),
		micro.Version("v1"),
		micro.WrapClient(opentracing.NewClientWrapper(t)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		// micro.Broker(brokerKafka.NewBroker(
		// 	broker.Option(func(opt *broker.Options) {
		// 		opt.Addrs = []string{"10.200.119.128:9092"}
		// 	}),
		// )),
	)

	service.Init(
		micro.RegisterTTL(1*time.Second),
		micro.RegisterInterval(1*time.Second),
	)

	return service
}
