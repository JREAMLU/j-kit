package elastic

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddr = "10.200.202.35:8500"
)

func TestLoadConfig(t *testing.T) {
	Convey("load elastic test", t, func() {
		err := Load(consulAddr, false, false, "BGCluster")
		So(err, ShouldBeNil)
		GetElastic("BGCluster")
	})
}
