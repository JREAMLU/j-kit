package mongo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddr   = "10.200.202.35:8500"
	instanceName = "BGBarrage"
)

func TestLoadConfig(t *testing.T) {
	Convey("load elastic test", t, func() {
		err := Load(consulAddr, false, instanceName)
		So(err, ShouldBeNil)
		t.Log(err)

		mgoClient := GetMongo(instanceName)
		So(mgoClient, ShouldNotBeNil)
		t.Log(mgoClient)
	})
}

func TestLoadAll(t *testing.T) {
	Convey("load elastic all test", t, func() {
		err := Load(consulAddr, false)
		So(err, ShouldBeNil)
		t.Log(err)

		mgoClients := GetAllMongo()
		So(len(mgoClients), ShouldBeGreaterThan, 0)
		t.Log(mgoClients)
		t.Log(mgoClients[instanceName])
	})
}
