package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringSplice(t *testing.T) {
	Convey("string splice test", t, func() {
		Convey("correct", func() {
			s := StringSplice("a", "b", "c", "d")
			So(s, ShouldEqual, "abcd")
		})
	})
}

func BenchmarkStringSplice(b *testing.B) {
	content := generateData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringSplice(content...)
	}
}

func generateData() []string {
	var s []string
	for i := 0; i < 100; i++ {
		s = append(s, "abc")
	}

	return s
}
