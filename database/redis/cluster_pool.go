package redis

import (
	"time"

	"github.com/mna/redisc"
)

var (
	clusterPool map[string]*redisc.Cluster
)

// @TODO pool
func getClusterPool(instanceName string, maxIdle int, idleTimeout time.Duration) (cluster *redisc.Cluster) {
	return nil
}
