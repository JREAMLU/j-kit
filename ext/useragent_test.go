package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseUserAgent(t *testing.T) {
	Convey("parse userAgent test", t, func() {
		Convey("correct", func() {
			m := ParseUserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
			So(m.IsBot, ShouldEqual, true)
			So(m.IsMobile, ShouldEqual, false)
			So(m.Mozilla, ShouldEqual, "5.0")
			So(m.Browser, ShouldEqual, "Googlebot")
			So(m.BrowserVersion, ShouldEqual, "2.1")
		})
	})
}
