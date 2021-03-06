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
	// ConfigNotExistsOrLoad redis config no exists
	ConfigNotExistsOrLoad = `redis config not exists OR not load config instance, server=%s,master=%v`
)

var settings map[string]*Group

var (
	// ConnectTimeout default redis Connect Timeout
	ConnectTimeout = 3 * time.Second
	// ReadTimeout default redis
	ReadTimeout = 1 * time.Second
	// WriteTimeout default redis
	WriteTimeout = 1 * time.Second
	// KeepAlivePeriod default keep alive period
	KeepAlivePeriod = 2 * time.Hour
)

func configNotExistsOrLoad(instanceName string, isMaster bool) error {
	return fmt.Errorf(ConfigNotExistsOrLoad, instanceName, isMaster)
}

// SetConnectTimeout Set Connect Timeout
func SetConnectTimeout(t time.Duration) {
	ConnectTimeout = t
}

// SetKeepAlivePeriod Set Keep Alive Period
func SetKeepAlivePeriod(t time.Duration) {
	KeepAlivePeriod = t
}

// SetReadTimeout Set Read Timeout
func SetReadTimeout(t time.Duration) {
	ReadTimeout = t
}

// SetWriteTimeout Set Write Timeout
func SetWriteTimeout(t time.Duration) {
	WriteTimeout = t
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

func getClusterNodes(instanceName string) []string {
	if _, ok := settings[instanceName]; !ok {
		return nil
	}

	nodes := make([]string, len(settings[instanceName].RedisConns))
	for k, v := range settings[instanceName].RedisConns {
		nodes[k] = v.ConnStr
	}

	return nodes
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

		LB := loadBalance(len(pool))
		hited := pool[LB]
		return &group.RedisConns[hited]
	}

	return nil
}

func loadBalance(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}
