package consul

import (
	"path"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/constant"
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

// KafkaBrokers brokers
type KafkaBrokers struct {
	Brokers []string
}

// KafkaZookeeper zookeeper zkroot
type KafkaZookeeper struct {
	Addrs  []string
	Zkroot string
}

// GetKafkas get kafka addrs
func (client *Client) GetKafkas(clusterName string) ([]string, error) {
	buf, err := client.Get(path.Join(Kafka, clusterName))
	var brokers KafkaBrokers
	if err != nil {
		return nil, err
	}

	_, err = toml.Decode(buf, &brokers)
	if err != nil {
		return nil, err
	}

	return brokers.Brokers, nil
}

// GetZookeepers get zookeeper addrs
func (client *Client) GetZookeepers(clusterName string) ([]string, string, error) {
	buf, err := client.Get(path.Join(Zookeeper, clusterName))
	var zk KafkaZookeeper
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	_, err = toml.Decode(buf, &zk)
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	return zk.Addrs, zk.Zkroot, nil
}
