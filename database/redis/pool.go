package redis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pools map[string]*Pool
var rwMutex sync.RWMutex

// Pool redis pool
type Pool struct {
	pools       map[string]*redis.Pool
	addr        string
	MaxIdle     int
	IdleTimeout time.Duration
	rwMutex     sync.RWMutex
}

func newPool(addr string, maxIdle int, idleTimeout time.Duration) *Pool {
	return &Pool{
		pools:       make(map[string]*redis.Pool),
		addr:        addr,
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
	}
}

func (p *Pool) get(db string) (pool *redis.Pool) {
	p.rwMutex.RLock()
	pool = p.pools[db]
	p.rwMutex.RUnlock()
	return
}

// GetPool get redis pool
func GetPool(addr, db string, maxIdle int, idleTimeout time.Duration) *redis.Pool {
	// @TODO pool
	return nil
}
