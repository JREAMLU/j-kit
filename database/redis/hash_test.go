package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	hashServer       = "test"
	hashHost         = "10.200.202.35:8500"
	hashKeyPrefixFmt = "%v"
	hashKey          = "thash"
)

func TestGet(t *testing.T) {
	Convey("get test", t, func() {
		h := NewHash(hashServer, hashKeyPrefixFmt)
		reply, err := h.Get(hashKey, "name")
		t.Log(reply)
		t.Log(err)
	})
}
