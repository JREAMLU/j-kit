package consul

import (
	"log"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/constant"
)

// Conn prefix conn
const (
	Conn           = "conn"
	SeparatorStart = `<<<--------------------------------------------------------`
	SeparatorEnd   = `-------------------------------------------------------->>>`
)

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
	// Consul mongo connect
	Consul = path.Join(Conn, "v1/consul")
)

// KafkaBrokers brokers
type KafkaBrokers struct {
	Brokers []string
	Topic   string
}

// KafkaZookeeper zookeeper zkroot
type KafkaZookeeper struct {
	Addrs  []string
	Zkroot string
}

// RegistryConsul registry consul
type RegistryConsul struct {
	Addrs []string
}

// GetKafkas get kafka addrs
func (client *Client) GetKafkas(clusterName string) ([]string, string, error) {
	key := path.Join(Kafka, clusterName)
	buf, err := client.Get(key)
	var brokers KafkaBrokers
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	log.Printf("Load Kafka Config: %v \n%v\n%v\n%v\n\n", key, SeparatorStart, buf, SeparatorEnd)

	_, err = toml.Decode(buf, &brokers)
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	return brokers.Brokers, brokers.Topic, nil
}

// GetZookeepers get zookeeper addrs
func (client *Client) GetZookeepers(clusterName string) ([]string, string, error) {
	key := path.Join(Zookeeper, clusterName)
	buf, err := client.Get(path.Join(Zookeeper, clusterName))
	var zk KafkaZookeeper
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	log.Printf("Load Zookeeper Config: %v \n%v\n%v\n%v\n\n", key, SeparatorStart, buf, SeparatorEnd)

	_, err = toml.Decode(buf, &zk)
	if err != nil {
		return nil, constant.EmptyStr, err
	}

	return zk.Addrs, zk.Zkroot, nil
}

// GetConsulAddrs get consul addrs
func (client *Client) GetConsulAddrs(clusterName string) ([]string, error) {
	key := path.Join(Consul, clusterName)
	buf, err := client.Get(key)
	var rc RegistryConsul
	if err != nil {
		return nil, err
	}

	log.Printf("Load Consul Config: %v \n%v\n%v\n%v\n\n", key, SeparatorStart, buf, SeparatorEnd)

	_, err = toml.Decode(buf, &rc)
	if err != nil {
		return nil, err
	}

	return rc.Addrs, nil
}
