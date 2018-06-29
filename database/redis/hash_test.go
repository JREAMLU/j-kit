package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddrHash   = "10.200.202.35:8500"
	hashServer       = "Crawler"
	hashKeyPrefixFmt = "%v"
	hashKey          = "thash"
)

func TestGet(t *testing.T) {
	Convey("get test", t, func() {
		Load(consulAddrHash, false)

		h := NewHash(hashServer, hashKeyPrefixFmt)
		reply, err := h.Get(hashKey, "name")
		t.Log(reply)
		t.Log(err)
	})
}
