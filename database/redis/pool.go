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

func (p *Pool) get(db string) (rPool *redis.Pool) {
	p.rwMutex.RLock()
	rPool = p.pools[db]
	p.rwMutex.RUnlock()

	return rPool
}

// Get get or register pool
func (p *Pool) Get(db string) (rPool *redis.Pool) {
	rPool = p.get(db)
	if rPool != nil {
		return rPool
	}

	return p.register(db)
}

// @TODO register pool
func (p *Pool) register(db string) (rPool *redis.Pool) {
	return nil
}

// GetPool get redis pool Ingress
func GetPool(addr, db string, maxIdle int, idleTimeout time.Duration) *redis.Pool {
	pool := getPool(addr, db)
	if pool != nil {
		return pool
	}

	return registerPool(addr, maxIdle, idleTimeout).Get(db)
}

func getPool(addr, db string) (rPool *redis.Pool) {
	rwMutex.RLock()
	if pool, ok := pools[addr]; ok {
		rPool = pool.Get(db)
	}
	rwMutex.RUnlock()

	return rPool
}

func registerPool(addr string, maxIdle int, idleTimeout time.Duration) *Pool {
	var (
		pool *Pool
		ok   bool
	)

	rwMutex.Lock()
	if pool, ok = pools[addr]; !ok {
		pool = newPool(addr, maxIdle, idleTimeout)
		pools[addr] = pool
	}
	rwMutex.Unlock()

	return pool
}
