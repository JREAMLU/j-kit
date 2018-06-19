package consul

import (
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

// WatchKey watch key
func WatchKey(consulAddr, key string, handle func(*api.KVPair)) {
	plan, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  key,
	})

	if err != nil {
		panic(err)
	}

	first := true
	plan.Handler = func(idx uint64, raw interface{}) {
		// when first run, key exist and do handler
		if first {
			first = false
			return
		}

		if raw == nil {
			return
		}

		v, ok := raw.(*api.KVPair)
		if ok && v != nil {
			handle(v)
		}
	}

	err = plan.Run(consulAddr)
	if err != nil {
		panic(err)
	}
}

// WatchKeyPrefix watch dir
func WatchKeyPrefix(consulAddr, key string, handle func(api.KVPairs)) {
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}

	plan, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": key,
	})
	if err != nil {
		panic(err)
	}

	first := true
	plan.Handler = func(idx uint64, raw interface{}) {
		if first {
			first = false
			return
		}

		if raw == nil {
			return
		}

		v, ok := raw.(api.KVPairs)
		if ok && len(v) > 0 {
			handle(v)
		}
	}

	err = plan.Run(consulAddr)
	if err != nil {
		panic(err)
	}
}
