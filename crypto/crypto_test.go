package crypto

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMD5(t *testing.T) {
	Convey("md5 test", t, func() {
		Convey("lower", func() {
			encrypt, err := MD5("123")
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "202cb962ac59075b964b07152d234b70")
		})

		Convey("upper", func() {
			encrypt, err := MD5("123", true)
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "202CB962AC59075B964B07152D234B70")
		})
	})

	Convey("sha1 test", t, func() {
		Convey("lower", func() {
			encrypt, err := Sha1("123")
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "40bd001563085fc35165329ea1ff5c5ecbdbbeef")
		})

		Convey("upper", func() {
			encrypt, err := Sha1("123", true)
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "40BD001563085FC35165329EA1FF5C5ECBDBBEEF")
		})
	})

	Convey("hmac-md5 test", t, func() {
		Convey("lower", func() {
			encrypt, err := HMacMD5([]byte("123"), "abc")
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "725658455c63977e1b73a199970a9972")
		})

		Convey("upper", func() {
			encrypt, err := HMacMD5([]byte("123"), "abc", true)
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "725658455C63977E1B73A199970A9972")
		})
	})

	Convey("hmac-sha1 test", t, func() {
		Convey("lower", func() {
			encrypt, err := HMacSha1([]byte("123"), "abc")
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "be9106a650ede01f4a31fde2381d06f5fb73e612")
		})

		Convey("upper", func() {
			encrypt, err := HMacSha1([]byte("123"), "abc", true)
			So(err, ShouldBeNil)
			So(encrypt, ShouldEqual, "BE9106A650EDE01F4A31FDE2381D06F5FB73E612")
		})
	})
}
