package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
)

func traceIntoContext(ctx context.Context, tracer opentracing.Tracer, name string, req *http.Request) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	var span opentracing.Span
	wireContext, err := tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span = tracer.StartSpan(name)
	} else {
		span = tracer.StartSpan(name, opentracing.ChildOf(wireContext))
	}

	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return nil, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)
	return ctx, span, nil
}

func traceIntoContextCall(ctx context.Context, tracer opentracing.Tracer, name string, req *http.Request) (context.Context, opentracing.Span, error) {
	span := opentracing.SpanFromContext(req.Context())
	span = span.Tracer().StartSpan(name, opentracing.ChildOf(span.Context()))

	// Inject the Span context into the outgoing HTTP Request.
	if err := tracer.Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	); err != nil {
		return nil, nil, err
	}

	return ctx, span, nil
}

// RequestFunc is a middleware function for outgoing HTTP requests.
type RequestFunc func(req *http.Request) *http.Request

// CallHTTPRequest to http
func CallHTTPRequest(tracer opentracing.Tracer) RequestFunc {
	return func(req *http.Request) *http.Request {
		name := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)

		ctx, span, err := traceIntoContextCall(req.Context(), tracer, name, req)
		if err != nil {
			return req
		}
		defer span.Finish()
		span.LogEvent("call")

		return req.WithContext(ctx)
	}
}

// HandlerFunc is a middleware function for incoming HTTP requests.
type HandlerFunc func(next http.Handler) http.Handler

// HandlerHTTPRequest req
func HandlerHTTPRequest(tracer opentracing.Tracer, operationName string) HandlerFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx, span, err := traceIntoContext(req.Context(), tracer, operationName, req)
			if err != nil {
				return
			}
			defer span.Finish()

			span.LogEvent("handler")

			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}

// HandlerHTTPRequestGin gin
func HandlerHTTPRequestGin(tracer opentracing.Tracer, operationName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span, err := traceIntoContext(c.Request.Context(), tracer, operationName, c.Request)
		if err != nil {
			return
		}
		defer span.Finish()

		span.LogEvent("handler")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
