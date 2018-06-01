package useragent

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var ual = []string{
	` `,
	`123`,
	`/`,
	`Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36 QIHU 360EE`,
	`Mozilla/5.0 (Linux; Android 5.1; 1501_M02 Build/LMY47D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.94 Mobile Safari/537.36 360 Aphone Browser (100.2.0)`,
	`Mozilla/5.0 (iPhone 5ATT; CPU iPhone OS 9_3_2 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/6.0 MQQBrowser/6.8.1 Mobile/13F69 Safari/8536.25 MttCustomUA/2`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko LBBROWSER`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.22 Safari/537.36 SE 2.X MetaSr 1.0`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.108 Safari/537.36 2345Explorer/8.0.0.13547`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.87 Safari/537.36 QQBrowser/9.2.5748.400`,
	`Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 UBrowser/5.7.15319.5 Safari/537.36`,
	`Mozilla/5.0 (Linux; U; Android 4.4.4; zh-CN; vivo Y13L Build/KTU84P) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 UCBrowser/11.0.2.840 U3/0.8.0 Mobile Safari/534.30`,
	`Mozilla/5.0 (Linux; Android 6.0; Le X620 Build/HEXCNFN5801708221S; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/49.0.2623.91 Mobile Safari/537.36 AliApp(CN/3.6.2) WindVane/8.0.0`,
	`Mozilla/5.0 (Linux; Android 4.4.2; H30-T00 Build/HuaweiH30-T00) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 AliApp(TB/4.9.2) WindVane/5.2.2 TBANDROID/700342@taobao_android_4.9.2`,
}

func TestLoadRedisConfig(t *testing.T) {
	Convey("func LoadRedisConfig()", t, func() {
		Convey("correct", func() {
			for _, ua := range ual {
				So(func() { New(ua) }, ShouldNotPanic)
			}
		})
	})
}

func TestParseUA(t *testing.T) {
	Convey("func ParseUA()", t, func() {
		Convey("UBrowser", func() {
			o := New(ual[10])
			bn, _ := o.Browser()
			So(bn, ShouldEqual, "UC")
		})
		Convey("UCBrowser", func() {
			o := New(ual[11])
			bn, _ := o.Browser()
			So(bn, ShouldEqual, "UC")
		})
		Convey("WindVane", func() {
			o := New(ual[12])
			bn, _ := o.Browser()
			So(bn, ShouldEqual, "AliApp")
		})
		Convey("AliApp", func() {
			o := New(ual[13])
			bn, _ := o.Browser()
			So(bn, ShouldEqual, "AliApp")
		})
	})
}
