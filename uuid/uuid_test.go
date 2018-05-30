package uuid

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerate(t *testing.T) {
	Convey("Generate test", t, func() {
		Convey("correct", func() {
			uuid, err := Generate()
			So(err, ShouldBeNil)
			So(uuid, ShouldNotBeEmpty)
		})
	})
}
