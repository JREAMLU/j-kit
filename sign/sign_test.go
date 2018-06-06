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
	tt := time.Now().Unix()
	tt = 1467012356
	expire := time.Now().Unix() + int64(1000)
	expired := int64(1000)

	Convey("sign test", t, func() {
		Convey("Generate Sign", func() {
			var err error
			signed, err = GenerateSign([]byte(data), tt, secret)
			So(err, ShouldBeNil)
			So(signed, ShouldNotBeNil)
			So(signed, ShouldEqual, "76E152FF759112893A160C71E8646726")
		})

		Convey("Valid Sign correct", func() {
			err := ValidSign([]byte(data), signed, tt, secret, expire)
			So(err, ShouldBeNil)
		})

		Convey("Valid Sign incorrect", func() {
			err := ValidSign([]byte(data), "xxoo", tt, secret, expire)
			So(err, ShouldNotBeNil)
		})

		Convey("Valid Sign expired", func() {
			err := ValidSign([]byte(data), signed, tt, secret, expired)
			So(err, ShouldNotBeNil)
		})
	})
}
