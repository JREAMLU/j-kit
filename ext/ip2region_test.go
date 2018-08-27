package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuery(t *testing.T) {
	region, err := NewRegion("ip2region.db")
	if err != nil {
		panic(err)
	}

	Convey("ip query", t, func() {
		ip, err := region.Query([]string{"127.0.0.1", "119.75.218.70"}, "memory")
		So(err, ShouldBeNil)
		So(ip, ShouldNotBeNil)
	})
}
