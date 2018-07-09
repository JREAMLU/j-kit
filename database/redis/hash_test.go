package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddrHash    = "10.200.202.35:8500"
	hashServer        = "Crawler"
	hashServerCluster = "CrawlerCluster"
	hashKeyPrefixFmt  = "%v"
	hashKey           = "thash"
)

func TestGet(t *testing.T) {
	// Convey("get test", t, func() {
	// 	Load(consulAddrHash, false, hashServer)
	//
	// 	h := NewHash(hashServer, hashKeyPrefixFmt)
	// 	reply, err := h.Get(hashKey, "name")
	// 	So(err, ShouldBeNil)
	// 	So(reply, ShouldNotBeBlank)
	// 	t.Log(reply, err)
	// })

	// Convey("get cluster test", t, func() {
	// 	Load(consulAddrHash, false, hashServerCluster)
	//
	// 	h := NewHash(hashServerCluster, hashKeyPrefixFmt)
	// 	reply, err := h.Get(hashKey, "name")
	// 	So(err, ShouldBeNil)
	// 	So(reply, ShouldNotBeBlank)
	// 	t.Log(reply, err)
	// })

	Convey("gets test", t, func() {
		Load(consulAddrHash, false, hashServer)

		h := NewHash(hashServer, hashKeyPrefixFmt)
		// var fields []string{"name","age"}
		fields := []string{"name", "age"}
		reply, err := h.Gets(hashKey, fields)
		So(err, ShouldBeNil)
		So(len(reply), ShouldBeGreaterThan, 0)
		t.Log(reply, err)
	})
}
