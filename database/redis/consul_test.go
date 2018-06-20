package redis

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddr = "10.200.202.35:8500"
)

func TestLoadConfig(t *testing.T) {
	Convey("load redis test", t, func() {
		Convey("load all", func() {
			err := Load(consulAddr, false)
			fmt.Println("++++++++++++: 1", err)
		})
	})
}
