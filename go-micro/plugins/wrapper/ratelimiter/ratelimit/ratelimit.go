package ratelimit

import (
	"time"

	"github.com/juju/ratelimit"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"

	"context"
)

type clientWrapper struct {
	fn func(req server.Request) error
	client.Client
}

func clientLimit(bs map[string]*ratelimit.Bucket, wait bool, errID string) func(req server.Request) error {
	return func(req server.Request) error {
		if _, ok := bs[req.Service()]; ok {
			if wait {
				time.Sleep(bs[req.Service()].Take(1))
			} else if bs[req.Service()].TakeAvailable(1) == 0 {
				return errors.New(errID, "too many request", 429)
			}
		}

		return nil
	}
}

func handlerLimit(b *ratelimit.Bucket, wait bool, errID string) func() error {
	return func() error {
		if wait {
			time.Sleep(b.Take(1))
		} else if b.TakeAvailable(1) == 0 {
			return errors.New(errID, "too many request", 429)
		}

		return nil
	}
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if err := c.fn(req); err != nil {
		return err
	}
	return c.Client.Call(ctx, req, rsp, opts...)
}

// NewClientWrapper takes a rate limiter and wait flag and returns a client Wrapper.
func NewClientWrapper(bs map[string]*ratelimit.Bucket, wait bool) client.Wrapper {
	fn := clientLimit(bs, wait, "go.micro.client")

	return func(c client.Client) client.Client {
		return &clientWrapper{fn, c}
	}
}

// NewHandlerWrapper takes a rate limiter and wait flag and returns a client Wrapper.
func NewHandlerWrapper(b *ratelimit.Bucket, wait bool) server.HandlerWrapper {
	fn := handlerLimit(b, wait, "go.micro.server")

	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if err := fn(); err != nil {
				return err
			}
			return h(ctx, req, rsp)
		}
	}
}
