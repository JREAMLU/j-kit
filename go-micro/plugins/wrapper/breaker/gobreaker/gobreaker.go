package gobreaker

import (
	"github.com/micro/go-micro/client"
	"github.com/sony/gobreaker"

	"context"
)

type clientWrapper struct {
	cbs map[string]*gobreaker.CircuitBreaker
	client.Client
}

// Call call
// request servcie exist, do circuitBreaker
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if _, ok := c.cbs[req.Service()]; ok {
		_, err := c.cbs[req.Service()].Execute(func() (interface{}, error) {
			cerr := c.Client.Call(ctx, req, rsp, opts...)
			return nil, cerr
		})

		return err
	}

	err := c.Client.Call(ctx, req, rsp, opts...)
	return err
}

// NewClientWrapper takes a *gobreaker.CircuitBreaker and returns a client Wrapper.
func NewClientWrapper(cbs map[string]*gobreaker.CircuitBreaker) client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{cbs, c}
	}
}
