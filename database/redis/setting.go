package redis

import "fmt"

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

func disableRefreshPool(instanceName string) {
	if group, ok := settings[instanceName]; ok {
		group.RefreshPool = false
	}
}
