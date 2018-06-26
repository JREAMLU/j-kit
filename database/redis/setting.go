package redis

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

var settings map[string]*Group
