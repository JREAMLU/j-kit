package http

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRollingCurl(t *testing.T) {
	Convey("func RollingCurl()", t, func() {
		Convey("correct", func() {
			res, err := request(15)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})

		Convey("uncorrect", func() {
			res, err := request(1)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeEmpty)
		})
	})
}

func request(timeout int64) (Responses, error) {
	req := NewRequests()
	resp, err := req.RequestCURL(
		"POST",
		"http://localhost/study/curl/servera.php",
		map[string]string{
			"Content-Type": "application/json;charset=UTF-8;",
		},
		`{"name":"KII","age":24}`,
		3,
		nil,
	)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
