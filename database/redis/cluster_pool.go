package redis

import (
	"net"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
)

var rwMutexC sync.RWMutex
var poolcs map[string]*redisc.Cluster

// Poolc cluster pool
type Poolc struct {
	rwMutex     sync.RWMutex
	MaxIdle     int
	IdleTimeout time.Duration
}

func init() {
	poolcs = make(map[string]*redisc.Cluster)
}

func newPoolc(maxIdle int, idleTimeout time.Duration) *Poolc {
	return &Poolc{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
	}
}

func getPoolc(instanceName string, maxIdle int, idleTimeout time.Duration) (cluster *redisc.Cluster) {
	rwMutexC.RLock()
	if c, ok := poolcs[instanceName]; ok {
		cluster = c
		rwMutexC.RUnlock()
		return cluster
	}
	rwMutexC.RUnlock()

	return registerPoolc(instanceName, maxIdle, idleTimeout)
}

func registerPoolc(instanceName string, maxIdle int, idleTimeout time.Duration) (cluster *redisc.Cluster) {
	nodes := getClusterNodes(instanceName)
	if len(nodes) == 0 {
		return nil
	}

	rwMutexC.RLock()
	poolc := newPoolc(maxIdle, idleTimeout)
	cluster = &redisc.Cluster{
		StartupNodes: nodes,
		DialOptions: []redis.DialOption{
			redis.DialConnectTimeout(ConnectTimeout),
			redis.DialReadTimeout(ReadTimeout),
			redis.DialWriteTimeout(WriteTimeout),
		},
		CreatePool: poolc.register,
	}
	rwMutexC.RUnlock()

	return cluster
}

func (p *Poolc) register(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     p.MaxIdle,
		MaxActive:   p.MaxIdle,
		IdleTimeout: p.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			d := net.Dialer{
				Timeout: ConnectTimeout,
			}
			dialConn, err := d.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			tcpConn := dialConn.(*net.TCPConn)
			if err = tcpConn.SetKeepAlive(true); err != nil {
				return nil, err
			}
			if err = tcpConn.SetKeepAlivePeriod(KeepAlivePeriod); err != nil {
				return nil, err
			}

			return redis.NewConn(tcpConn, ReadTimeout, WriteTimeout), nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}, nil
}
