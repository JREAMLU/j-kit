package redis

import (
	"fmt"
	"math/rand"
	"time"
)

// Group redis group
type Group struct {
	Name           string
	PoolSize       int64
	RedisConns     []Conn
	IsCluster      bool
	RefreshSetting bool
	RefreshPool    bool
}

// Conn redis conn
type Conn struct {
	ConnStr  string
	DB       string
	IsMaster bool
}

const (
	// ConfigNotExists redis config no exists
	ConfigNotExists = `redis config not exists,server=%s,master=%v`
)

var settings map[string]*Group

func configNotExists(instanceName string, isMaster bool) error {
	return fmt.Errorf(ConfigNotExists, instanceName, isMaster)
}

func isCluster(instanceName string) bool {
	if group, ok := settings[instanceName]; ok {
		return group.IsCluster
	}

	return false
}

func isRefreshPool(instanceName string) bool {
	if group, ok := settings[instanceName]; ok {
		return group.RefreshPool
	}

	return false
}

func toggleRefreshPool(instanceName string, toggle bool) {
	if group, ok := settings[instanceName]; ok {
		group.RefreshPool = toggle
	}
}

func getConn(instanceName string, isMaster bool) *Conn {
	var pool []int
	if group, ok := settings[instanceName]; ok {
		for key := range group.RedisConns {
			if group.RedisConns[key].IsMaster == isMaster {
				pool = append(pool, key)
			}
		}

		if len(pool) == 0 {
			return nil
		}

		LB := poolBalance(len(pool))
		hited := pool[LB]
		return &group.RedisConns[hited]
	}

	return nil
}

func poolBalance(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}
