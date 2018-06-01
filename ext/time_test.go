package ext

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToday(t *testing.T) {
	Convey("today test", t, func() {
		Convey("correct", func() {
			t := Today()
			So(t, ShouldNotBeNil)
		})
	})
}

func TestCurrHour(t *testing.T) {
	Convey("currHourUnix test", t, func() {
		Convey("correct", func() {
			t := CurrHour()
			So(t, ShouldNotBeNil)
		})
	})
}

func TestFormat(t *testing.T) {
	Convey("format test", t, func() {
		Convey("yymmdd", func() {
			yymmdd := Formatyymmdd(time.Now())
			So(yymmdd, ShouldNotBeNil)
		})

		Convey("yyyymmdd", func() {
			yyyymmdd := Formatyyyymmdd(time.Now())
			So(yyyymmdd, ShouldNotBeNil)
		})

		Convey("yyyymm", func() {
			yyyymm := Formatyyyymm(time.Now())
			So(yyyymm, ShouldNotBeNil)
		})

		Convey("default", func() {
			d := FormatDefault(time.Now())
			So(d, ShouldNotBeNil)
		})
	})
}
