package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/JREAMLU/j-kit/ext"
	"github.com/JREAMLU/j-kit/go-micro/util"
	"github.com/JREAMLU/j-kit/uuid"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
)

// Requests struct
type Requests struct {
	HTTPClient   *http.Client
	TraceRequest RequestFunc
	Cb           *gobreaker.CircuitBreaker
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

	cbName                  = "http"
	cbMaxRequests    uint32 = 100
	cbInterval              = 30
	cbTimeout               = 90
	cbCountsRequests uint32 = 1000
	cbFailureRatio          = 0.6
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
		Cb:           defaultCircuitBreakerSetting(),
	}
}

func defaultCircuitBreakerSetting() *gobreaker.CircuitBreaker {
	gid, err := uuid.Generate()
	if err != nil {
		panic(err)
	}

	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        ext.StringSplice(cbName, "-", gid),
		MaxRequests: cbMaxRequests,
		Interval:    time.Duration(cbInterval) * time.Second,
		Timeout:     time.Duration(cbTimeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= cbCountsRequests && failureRatio >= cbFailureRatio
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
		},
	})
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

// CbRequestCURL CircuitBreaker curl
func (r *Requests) CbRequestCURL(ctx context.Context, Method string, URLStr string, Header map[string]string, Raw string, data interface{}) (rp Responses, err error) {
	res, err := r.Cb.Execute(func() (interface{}, error) {
		return r.RequestCURL(ctx, Method, URLStr, Header, Raw, data)
	})
	if err != nil {
		return rp, err
	}

	return res.(Responses), nil
}

//RequestCURL http url
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

	ctx = context.WithValue(context.WithValue(ctx, rawBody, Raw), headers, Header)
	req = r.TraceRequest(req.WithContext(ctx))

	for hkey, hval := range Header {
		req.Header.Set(hkey, hval)
	}

RELOAD:

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		util.TraceLog(req.Context(), fmt.Sprintf("REQUEST ERROR: %v, COUNTER: %v", err, i))
		i++
		if i < retryTimes {
			goto RELOAD
		}
		return rp, err
	}
	defer resp.Body.Close()
	rp.Response = resp

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
