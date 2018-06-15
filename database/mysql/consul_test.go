package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("load mysql test", t, func() {
		Convey("correct", func() {
			Load("10.200.202.35:8500", "BGCrawler")
		})
	})
}
