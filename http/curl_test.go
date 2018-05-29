package http

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestCURL(t *testing.T) {
	Convey("RequestCURL test", t, func() {
		Convey("no timeout", func() {
			res, err := request(5, 0, nil)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})

		Convey("timeout", func() {
			res, err := request(1, 0, nil)
			So(err, ShouldNotBeNil)
			So(res, ShouldNotBeEmpty)
		})

		Convey("retry", func() {
			res, err := request(1, 3, nil)
			So(err, ShouldNotBeNil)
			So(res, ShouldNotBeEmpty)
		})

		Convey("data is map not nil", func() {
			var jmap map[string]interface{}
			res, err := request(5, 0, &jmap)
			So(err, ShouldBeNil)
			So(res.Data, ShouldNotBeEmpty)
		})
	})
}

func request(timeout int64, retry int64, data interface{}) (Responses, error) {
	req := NewRequests(timeout)
	resp, err := req.RequestCURL(
		"POST",
		"http://localhost/study/curl/servera.php",
		map[string]string{
			"Content-Type": "application/json;charset=UTF-8;",
		},
		`{"name":"KII","age":24}`,
		retry,
		data,
	)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
