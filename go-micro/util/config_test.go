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
	Convey("new elastic test", t, func() {
		key := getServiceKey(serviceName, serviceVersion)
		var conf PusherConfig
		err := loadConfig(consulAddr, key, &conf)
		t.Log(err)
		t.Log(conf.Config)
		t.Log(conf.Pusher)
	})
}
