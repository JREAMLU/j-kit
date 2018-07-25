package http

import (
	"context"
	"net/http"

	"github.com/JREAMLU/j-kit/ext"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	// TargetSRV target srv
	TargetSRV = "TargetSRV"
	// FromSRV from srv
	FromSRV = "FromSRV"
	// Method method
	Method = "Method"
	// Proto proto
	Proto = "Proto"
	// RawBody rawbody
	RawBody = "RawBody"
	// Header header
	Header = "Header"
	// ContentType contenttype
	ContentType = "ContentType"
	// Params params
	Params = "Params"
	// Unknow unknow
	Unknow = "unknow"
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

	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, span, nil
}

// RequestFunc is a middleware function for outgoing HTTP requests.
type RequestFunc func(req *http.Request) *http.Request

// CallHTTPRequest to http
func CallHTTPRequest(tracer opentracing.Tracer) RequestFunc {
	return func(req *http.Request) *http.Request {
		url := ext.StringSplice(req.URL.Scheme, "://", req.URL.Host, req.URL.Path)
		ctx, span, err := traceIntoContextCall(req.Context(), tracer, url, req)
		if err != nil {
			return req
		}
		defer span.Finish()

		target := ext.StringSplice(req.URL.Scheme, "://", req.URL.Host, req.URL.RequestURI())

		// body
		span.LogKV(
			TargetSRV, target,
			Method, req.Method,
			Header, req.Context().Value(headers),
			RawBody, req.Context().Value(rawBody),
			Proto, req.Proto,
		)

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

		url := ext.StringSplice(c.Request.Host, c.Request.RequestURI)
		rawBody, err := c.GetRawData()
		if err != nil {
			return
		}

		span.LogKV(
			FromSRV, c.Request.RemoteAddr,
			TargetSRV, url,
			Method, c.Request.Method,
			Proto, c.Request.Proto,
			RawBody, string(rawBody),
			ContentType, c.GetHeader("Content-Type"),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
