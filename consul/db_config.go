package consul

import (
	"path"
)

// Conn prefix conn
const Conn = "conn"

var (
	// MYSQL mysql connect
	MYSQL = path.Join(Conn, "v1/mysql")
	// Kafka kafka connect
	Kafka = path.Join(Conn, "v1/kafka")
	// Zookeeper zk connect
	Zookeeper = path.Join(Conn, "v1/zookeeper")
	// ElasticSearch es connect
	ElasticSearch = path.Join(Conn, "v1/elasticsearch")
	// Redis redis connect
	Redis = path.Join(Conn, "v1/redis")
	// MongoDB mongo connect
	MongoDB = path.Join(Conn, "v1/mongodb")
)
