package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIAbs(t *testing.T) {
	Convey("iabs int test", t, func() {
		Convey("correct", func() {
			i := IAbs(-3)
			So(i, ShouldEqual, 3)
		})
	})

	Convey("iabs int32 test", t, func() {
		Convey("correct", func() {
			i := IAbs32(-3)
			So(i, ShouldEqual, 3)
		})
	})

	Convey("iabs int64 test", t, func() {
		Convey("correct", func() {
			i := IAbs64(-3)
			So(i, ShouldEqual, 3)
		})
	})
}
