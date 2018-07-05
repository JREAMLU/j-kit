package redis

import (
	"net"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mna/redisc"
)

var (
	clusterPool        map[string]*redisc.Cluster
	clusterRwMutex     sync.RWMutex
	clusterMaxIdle     int
	clusterIdleTimeout time.Duration
)

// @TODO pool
func getClusterPool(instanceName string, maxIdle int, idleTimeout time.Duration) (cluster *redisc.Cluster) {
	return nil
}

func registerClusterPool(instanceName string, maxIdle int, idleTimeout time.Duration) (cluster *redisc.Cluster) {
	return nil
}

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     clusterMaxIdle,
		MaxActive:   clusterMaxIdle,
		IdleTimeout: clusterIdleTimeout,
		Dial: func() (redis.Conn, error) {
			/*
			   addr, err := net.ResolveTCPAddr("tcp", p.addr)
			   if err != nil {
			       return nil, err
			   }
			   tcpConn, err := net.DialTCP("tcp", nil, addr)
			*/
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
