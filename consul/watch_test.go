package consul

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testConsulAddrWatch = "10.200.202.35:8500"
	testKey             = "conn/v1/mysql/BGColllector"
	testKeyPrefix       = "conn/v1/mysql"
	testKeyWatch        = "conn/v1/mysql/Abc"
	testKeyPrefixWatch  = "conn/v1/mysql"
	testValueWatch      = "i am test value watch"
)

var (
	testClientWatch *Client
)

func init() {
	var err error
	if testClientWatch, err = NewClient(SetAddress(testConsulAddrWatch)); err != nil {
		panic(err)
	}
}

func TestWatch(t *testing.T) {
	Convey("watch test", t, func() {
		Convey("watch key", func() {
			err := testWatch(t, testKeyWatch, func(wg *sync.WaitGroup) {
				WatchKey(testConsulAddrWatch, testKeyWatch, func(kvPair *api.KVPair) {
					t.Log(kvPair.Key)
					t.Log(string(kvPair.Value))
					wg.Done()
				})
			})

			So(err, ShouldBeNil)
		})

		Convey("watch keyprefix", func() {
			err := testWatch(t, filepath.Join(testKeyPrefixWatch, "ip"), func(wg *sync.WaitGroup) {
				WatchKeyPrefix(testConsulAddrWatch, testKeyPrefixWatch, func(kvPairs api.KVPairs) {
					for k := range kvPairs {
						t.Log(kvPairs[k].Key)
						t.Log(string(kvPairs[k].Value))
					}
					wg.Done()
				})
			})

			So(err, ShouldBeNil)
		})
	})
}

func testWatch(t *testing.T, testKeyWatch string, f func(*sync.WaitGroup)) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go f(wg)
	time.Sleep(3 * time.Millisecond)
	if err := testClientWatch.Put(testKeyWatch, testValueWatch); err != nil {
		return err
	}
	wg.Wait()

	return nil
}
