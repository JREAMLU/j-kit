package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/JREAMLU/j-kit/ext"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
)

const (
	_defaultPagesize          = 500
	_defaultMaxidle           = 50
	_defaultIdletimeout       = 240 * time.Second
	_defaultClusterRetryTime  = 3
	_defaultClusterRetryDelay = 100 * time.Millisecond
)

// Structure redis structure
type Structure struct {
	KeyPrefixFmt string
	InstanceName string
	readPool     *redis.Pool
	writePool    *redis.Pool
	clusterPool  *redisc.Cluster
	writeConn    string
	readConn     string
	mutex        sync.Mutex
	MaxIdle      int
	IdleTimeout  time.Duration
}

// NewStructure new structure
func NewStructure(instanceName, keyPrefixFmt string) Structure {
	return Structure{
		KeyPrefixFmt: keyPrefixFmt,
		InstanceName: instanceName,
		MaxIdle:      _defaultMaxidle,
		IdleTimeout:  _defaultIdletimeout,
	}
}

// SetMaxIdle set max idle
func (s *Structure) SetMaxIdle(maxIdle int) {
	s.MaxIdle = maxIdle
}

// SetIdleTimeout set idle timeout
func (s *Structure) SetIdleTimeout(idleTimeout time.Duration) {
	s.IdleTimeout = idleTimeout
}

// InitKey init redis key
func (s *Structure) InitKey(keySuffix string) string {
	if ext.StringEq(keySuffix) {
		return s.KeyPrefixFmt
	}

	return fmt.Sprintf(s.KeyPrefixFmt, keySuffix)
}

func (s *Structure) String(isMaster bool, cmd string, params ...interface{}) (reply string, err error) {
	conn := s.getConn(isMaster)
	if conn == nil {
		return "", configNotExists(s.InstanceName, isMaster)
	}

	reply, err = redis.String(conn.Do(cmd, params...))
	conn.Close()
	return
}

func (s *Structure) getConn(isMaster bool) redis.Conn {
	if s.isCluster() {
		return s.getClusterConn()
	}

	return s.getClientConn(isMaster)
}

func (s *Structure) isCluster() bool {
	return isCluster(s.InstanceName)
}

func (s *Structure) getClientConn(isMaster bool) redis.Conn {
	// refresh true, set pool = nil, then get new pool
	if isRefreshPool(s.InstanceName) {
		s.mutex.Lock()
		s.writePool = nil
		s.readPool = nil
		toggleRefreshPool(s.InstanceName, false)
		s.mutex.Unlock()
	}

	if s.writePool == nil {
		s.writePool = s.getPool(s.InstanceName, true)
		s.readPool = s.getPool(s.InstanceName, false)
	}

	if isMaster {
		if s.writePool == nil {
			return nil
		}

		return s.writePool.Get()
	}

	if s.readPool == nil {
		return nil
	}

	return s.readPool.Get()
}

func (s *Structure) getClusterConn() redis.Conn {
	// refresh true, set pool = nil, then get new pool
	if isRefreshPool(s.InstanceName) {
		s.mutex.Lock()
		if s.clusterPool != nil {
			s.clusterPool.Close()
			delete(clusterPool, s.InstanceName)
			s.clusterPool = nil
		}
		toggleRefreshPool(s.InstanceName, false)
		s.mutex.Unlock()
	}

	if s.clusterPool == nil {
		s.clusterPool = s.getClusterPool(s.InstanceName)
	}

	retryConn, err := redisc.RetryConn(s.clusterPool.Get(), _defaultClusterRetryTime, _defaultClusterRetryDelay)
	if err != nil {
		return nil
	}

	return retryConn
}

func (s *Structure) getPool(instanceName string, isMaster bool) *redis.Pool {
	conn := getConn(instanceName, isMaster)
	if conn == nil {
		return nil
	}

	return GetPool(conn.ConnStr, conn.DB, s.MaxIdle, s.IdleTimeout)
}

func (s *Structure) getClusterPool(instanceName string) *redisc.Cluster {
	return getClusterPool(instanceName, s.MaxIdle, s.IdleTimeout)
}
