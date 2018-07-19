package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIP2Int(t *testing.T) {
	Convey("IP2Int test", t, func() {
		Convey("correct", func() {
			ip := "255.255.255.255"
			So(IP2Int(ip), ShouldEqual, 4294967295)
		})

		Convey("incorrect", func() {
			ip := "255.255.255.255.255"
			So(IP2Int(ip), ShouldEqual, 0)
		})
	})
}

func TestInt2IP(t *testing.T) {
	Convey("Int2IP test", t, func() {
		Convey("correct", func() {
			var ip int64 = 4294967295
			So(Int2IP(ip), ShouldEqual, "255.255.255.255")
		})
	})
}

func TestServerIP(t *testing.T) {
	Convey("ServerIP test", t, func() {
		Convey("correct", func() {
			ip, err := ServerIP()
			So(err, ShouldBeNil)
			So(ip, ShouldNotBeEmpty)
		})
	})
}

func TestExtractIP(t *testing.T) {
	Convey("Extract IP test", t, func() {
		Convey("correct", func() {
			ip, err := ExtractIP("")
			So(err, ShouldBeNil)
			So(ip, ShouldNotBeEmpty)
		})
	})
}

func BenchmarkIP2Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ip := "255.255.255.255"
		ipInt := IP2Int(ip)
		Int2IP(ipInt)
	}
}

func BenchmarkServerIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ServerIP()
	}
}

func BenchmarkExtractAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIP("")
	}
}
