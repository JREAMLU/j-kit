package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	consulAddr     = "10.200.202.35:8500"
	serviceName    = "pusher"
	serviceVersion = "v1"
)

type PusherConfig struct {
	*Config

	Pusher struct {
		Auth   bool
		Secret string
	}
}

func TestLoadConfig(t *testing.T) {
	Convey("Load Config test", t, func() {
		config, err := LoadConfig(consulAddr, serviceName, serviceVersion)
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
		t.Log(err)
		t.Log(config)
	})
}

func TestLoadCustomConfig(t *testing.T) {
	Convey("Load Custom Config test", t, func() {
		var config PusherConfig
		err := LoadCustomConfig(consulAddr, serviceName, serviceVersion, &config)
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
		t.Log(err)
		t.Log(config.Config)
		t.Log(config.Pusher)
	})
}
