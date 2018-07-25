package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

// Requests struct
type Requests struct {
	HTTPClient   *http.Client
	TraceRequest RequestFunc
}

// Responses struct
type Responses struct {
	Response *http.Response
	Body     string
	Data     interface{}
}

var (
	// maxIdleConnections maxIdleConnections
	maxIdleConnections int64 = 100
	// maxConnectionIdleTime in pools, one connect can idle time
	maxConnectionIdleTime       = 60 * time.Second
	timeout               int64 = 3
	retryTimes            int64 = 1
)

// NewRequests new requests
func NewRequests(tracer opentracing.Tracer) *Requests {
	return &Requests{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: int(maxIdleConnections),
				IdleConnTimeout:     maxConnectionIdleTime,
			},
			Timeout: time.Duration(timeout) * time.Second,
		},
		TraceRequest: CallHTTPRequest(tracer),
	}
}

// SetTimeout http client timeout
func (r *Requests) SetTimeout(timeout int64) {
	r.HTTPClient.Timeout = time.Duration(timeout) * time.Second
}

// SetMaxIdleConnsPerHost set idle time
func (r *Requests) SetMaxIdleConnsPerHost(maxIdleConnections int64) {
	r.HTTPClient.Transport = &http.Transport{
		MaxIdleConnsPerHost: int(maxIdleConnections),
		IdleConnTimeout:     maxConnectionIdleTime,
	}
}

// SetIdleConnTimeout set idle timeout
func (r *Requests) SetIdleConnTimeout(maxConnectionIdleTime time.Duration) {
	r.HTTPClient.Transport = &http.Transport{
		MaxIdleConnsPerHost: int(maxIdleConnections),
		IdleConnTimeout:     maxConnectionIdleTime,
	}
}

// SetRetryTimes set retry times
func (r *Requests) SetRetryTimes(times int64) {
	retryTimes = times
}

type cstring string

const (
	rawBody cstring = "RawBody"
	headers cstring = "Header"
)

//RequestCURL http请求url
func (r *Requests) RequestCURL(ctx context.Context, Method string, URLStr string, Header map[string]string, Raw string, data interface{}) (rp Responses, err error) {
	var i int64
	req, err := http.NewRequest(
		Method,
		URLStr,
		strings.NewReader(Raw),
	)
	if err != nil {
		return rp, err
	}

	ctx = context.WithValue(ctx, rawBody, Raw)
	ctx = context.WithValue(ctx, headers, Header)
	req = r.TraceRequest(req.WithContext(ctx))

	for hkey, hval := range Header {
		req.Header.Set(hkey, hval)
	}

RELOAD:

	// TODO resp span
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		i++
		if i < retryTimes {
			goto RELOAD
		}
		return rp, err
	}
	rp.Response = resp

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp, err
	}

	rp.Body = string(body)

	if data != nil {
		err = json.Unmarshal(body, data)
		if err != nil {
			return rp, err
		}

		rp.Data = data
	}

	return rp, nil
}

// RequestRollingCURL batch curl
func (r *Requests) RequestRollingCURL(Method string, URLStr string, Header map[string]string, Raw string, RetryTimes int64, data interface{}) (rp Responses, err error) {
	return Responses{}, nil
}
