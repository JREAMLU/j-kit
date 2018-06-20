package consul

import (
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testClient *Client
)

const (
	testConsulAddr = "10.200.202.35:8500"
	mysqlToml      = `InstanceName = "BGCrawler"
DBName = "crawler"

[readwrite]
Server = "127.0.0.1"
Port = "3306"
UserID = "root"
Password = "123"
CharSet = "utf8mb4"

[readonly]
Server = "127.0.0.1"
Port = "3306"
UserID = "root"
Password = "123"
CharSet = "utf8mb4"`

	zookeeperToml = `[zookeeper]
Addrs = [
    "127.0.0.1:2181",
    "127.0.0.2:2181",
    "127.0.0.3:2181"
]

Zkroot = "/kafka"`

	redisToml = `[[master]]
Name = "crawler_cluster"
DB = "0"
IP = "127.0.0.1"
Port = "29010"
PoolSize = 10.0
IsCluster = true
IsMaster = true

[[master]]
Name = "crawler_cluster"
DB = "0"
IP = "127.0.0.1"
Port = "29011"
PoolSize = 10.0
IsCluster = true
IsMaster = true

[[slave]]
Name = "crawler_cluster"
DB = "0"
IP = "127.0.0.1"
Port = "29012"
PoolSize = 10.0
IsCluster = true
IsMaster = false

[[slave]]
Name = "crawler_cluster"
DB = "0"
IP = "127.0.0.1"
Port = "29013"
PoolSize = 10.0
IsCluster = true
IsMaster = false`

	kafkaToml = `[kafka]
Brokers = [
    "127.0.0.1:9092",
    "127.0.0.2:9092",
    "127.0.0.3:9092"
]

Zkaddrs = [
    "127.0.0.1:2181",
    "127.0.0.2:2181",
    "127.0.0.3:2181"
]

Zkroot = "/kafka"
Topic = "zlog"
Group = "zipkin"
CommitInterval = 10
ProcessingTimeout = 10
ZookeeperTimeout = 3000`
)

func init() {
	var err error
	if testClient, err = NewClient(SetAddress(testConsulAddr)); err != nil {
		panic(err)
	}

	initTestData()
}

func initTestData() {
	if err := testClient.Put(path.Join(MYSQL, "BGCrawler-Test"), mysqlToml); err != nil {
		panic(err)
	}

	if err := testClient.Put(path.Join(Zookeeper, "Zipkin-Test"), zookeeperToml); err != nil {
		panic(err)
	}

	if err := testClient.Put(path.Join(Redis, "BGCrawler-Cluster-Test"), redisToml); err != nil {
		panic(err)
	}

	if err := testClient.Put(path.Join(Kafka, "Zipkin-Test"), kafkaToml); err != nil {
		panic(err)
	}
}

func TestConsul(t *testing.T) {
	Convey("get test", t, func() {
		Convey("get mysql", func() {
			value, err := testClient.Get(path.Join(MYSQL, "BGCrawler-Test"))
			So(err, ShouldBeNil)
			So(value, ShouldNotBeEmpty)
		})

		Convey("get zookeeper", func() {
			value, err := testClient.Get(path.Join(Zookeeper, "Zipkin-Test"))
			So(err, ShouldBeNil)
			So(value, ShouldNotBeEmpty)
		})

		Convey("get redis", func() {
			value, err := testClient.Get(path.Join(Redis, "BGCrawler-Cluster-Test"))
			So(err, ShouldBeNil)
			So(value, ShouldNotBeEmpty)
		})

		Convey("get kafka", func() {
			value, err := testClient.Get(path.Join(Kafka, "Zipkin-Test"))
			So(err, ShouldBeNil)
			So(value, ShouldNotBeEmpty)
		})
	})

	Convey("delete test", t, func() {
		Convey("delete mysql", func() {
			err := testClient.Delete(path.Join(MYSQL, "BGCrawler-Test"))
			So(err, ShouldBeNil)
		})

		Convey("delete zookeeper", func() {
			err := testClient.Delete(path.Join(Zookeeper, "Zipkin-Test"))
			So(err, ShouldBeNil)
		})

		Convey("delete redis", func() {
			err := testClient.Delete(path.Join(Redis, "BGCrawler-Cluster-Test"))
			So(err, ShouldBeNil)
		})

		Convey("delete kafka", func() {
			err := testClient.Delete(path.Join(Kafka, "Zipkin-Test"))
			So(err, ShouldBeNil)
		})
	})
}
