package util

import (
	"context"
	"strconv"

	"github.com/JREAMLU/j-kit/ext"
	jopentracing "github.com/JREAMLU/j-kit/go-micro/trace/opentracing"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	_debug         = true
	_sameSpan      = true
	_traceID128Bit = true
)

// NewTrace new trace
func NewTrace(serviceName, version string, kafkaAddrs []string, kafkaTopic string, hostPort ...string) (opentracing.Tracer, error) {
	collector, err := zipkin.NewKafkaCollector(
		kafkaAddrs,
		zipkin.KafkaTopic(kafkaTopic),
	)
	if err != nil {
		return nil, err
	}

	var ipPort string
	if len(hostPort) == 0 {
		ipPort, err = ext.ExtractIP("")
	} else {
		ipPort, err = ext.ExtractIP(hostPort[0])
	}

	if err != nil {
		return nil, err
	}

	recorder := zipkin.NewRecorder(collector, _debug, ipPort, serviceName)
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(_sameSpan),
		zipkin.TraceID128Bit(_traceID128Bit),
	)
	if err != nil {
		return nil, err
	}

	opentracing.InitGlobalTracer(tracer)

	return tracer, nil
}

// TraceLog log
func TraceLog(ctx context.Context, logger string) {
	// toggle
	if md, ok := metadata.FromContext(ctx); ok {
		i, err := strconv.Atoi(md[jopentracing.ZipkinToggle])
		if err != nil || i <= 0 {
			return
		}
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	span.LogEvent(logger)
	span.Finish()
}

// TraceLogInject inject
func TraceLogInject(ctx context.Context, operationName string, logger string) {
	// toggle
	if md, ok := metadata.FromContext(ctx); ok {
		i, err := strconv.Atoi(md[jopentracing.ZipkinToggle])
		if err != nil || i <= 0 {
			return
		}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	if span == nil {
		return
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		return
	}

	span.LogEvent(ext.StringSplice("func(", operationName, ") ", logger))
	span.Finish()
}

// SetZipkin set zipkin trace
func SetZipkin(service micro.Service, kafkaAddrs []string, kafkaTopic string, hostPort ...string) {
	opts := service.Options()

	setZipkin(
		opts.Server.Options().Name,
		opts.Server.Options().Version,
		kafkaAddrs,
		kafkaTopic,
		hostPort...,
	)(&opts)
}

func setZipkin(serviceName, version string, kafkaAddrs []string, kafkaTopic string, hostPort ...string) micro.Option {
	return func(opt *micro.Options) {
		collector, err := zipkin.NewKafkaCollector(
			kafkaAddrs,
			zipkin.KafkaTopic(kafkaTopic),
		)
		if err != nil {
			panic(err)
		}

		var ipPort string
		if len(hostPort) == 0 {
			ipPort, err = ext.ExtractIP("")
		} else {
			ipPort, err = ext.ExtractIP(hostPort[0])
		}

		if err != nil {
			panic(err)
		}

		recorder := zipkin.NewRecorder(collector, _debug, ipPort, serviceName)
		tracer, err := zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(_sameSpan),
			zipkin.TraceID128Bit(_traceID128Bit),
		)
		if err != nil {
			panic(err)
		}

		clientWrap := jopentracing.NewClientWrapper(tracer)
		serverWrap := jopentracing.NewHandlerWrapper(tracer)

		micro.WrapClient(clientWrap)(opt)
		micro.WrapHandler(serverWrap)(opt)
	}
}
