package sign

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	secret = "abc123"
	data   = `{"name":"luj","age":13}`
)

func TestSign(t *testing.T) {
	var signed string

	Convey("sign test", t, func() {
		Convey("Generate Sign", func() {
			var err error
			signed, err = GenerateSign([]byte(data), time.Now().Unix(), secret)
			So(err, ShouldBeNil)
			So(signed, ShouldNotBeNil)
		})

		Convey("Valid Sign", func() {
			err := ValidSign([]byte(data), signed, time.Now().Unix(), secret, time.Now().Unix())
			So(err, ShouldBeNil)
		})
	})
}
