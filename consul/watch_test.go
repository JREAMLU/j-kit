package consul

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testConsulAddrWatch = "10.200.202.35:8500"
	testKey             = "conn/v1/mysql/BGColllector"
	testKeyPrefix       = "conn/v1/mysql"
)

func TestWatch(t *testing.T) {
	Convey("watch test", t, func() {
		Convey("watch key", func() {
			// 	WatchKey(testConsulAddrWatch, testKey, func(kvPair *api.KVPair) {
			// 		fmt.Println("++++++++++++: ", string(kvPair.Value))
			// 	})
		})

		Convey("watch keyprefix", func() {
			// WatchKeyPrefix(testConsulAddrWatch, testKeyPrefix, func(kvPairs api.KVPairs) {
			// 	for k := range kvPairs {
			// 		fmt.Println("++++++++++++: ", string(kvPairs[k].Value))
			// 	}
			// })
		})
	})
}
