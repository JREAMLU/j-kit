package redis

import (
	"fmt"
	"log"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/consul"
	"github.com/hashicorp/consul/api"
)

// MasterSlave master slave
type MasterSlave bool

// Configs redis configs in consul
type Configs struct {
	InstanceName string
	PoolSize     int64
	IsCluster    bool
	Master       []msConn
	Slave        []msConn
}

// msConn master & slave conn
type msConn struct {
	DB       string
	IP       string
	Port     string
	IsMaster bool
}

func (masterSlave MasterSlave) String() string {
	if masterSlave {
		return "master"
	}
	return "slave"
}

// Watch watch config
func Watch(consulAddr string, reloadConfig chan string, names ...string) {
	for i := range names {
		go func(name string) {
			consul.WatchKey(consulAddr, path.Join(consul.Redis, name), func(kvPair *api.KVPair) {
				reloadConfig <- name
			})
		}(names[i])
	}
}

func watching(consulAddr string, names ...string) {
	watchdNode := make(chan string)
	Watch(consulAddr, watchdNode, names...)
	go func() {
		for {
			select {
			case node := <-watchdNode:
				log.Printf("changed: %v \r\n", node)
				if err := LoadConfig(consulAddr, false, node); err != nil {
					log.Printf("Failed on redis LoadConfig, watchedNode: %v, err: %v \r\n", node, err)
					continue
				}

				if group, ok := settings[node]; ok {
					group.RefreshSetting = true
				}
			}
		}
	}()
}

// Load load redis
func Load(consulAddr string, isWatching bool, names ...string) error {
	return LoadConfig(consulAddr, isWatching, names...)
}

// LoadConfig load config
func LoadConfig(consulAddr string, isWatching bool, names ...string) error {
	client, err := consul.NewClient(consul.SetAddress(consulAddr))
	if err != nil {
		return err
	}

	if len(names) == 0 {
		if err = loadAll(client); err != nil {
			return err
		}
	} else {
		if err = loadByNames(client, names); err != nil {
			return err
		}
	}

	if isWatching {
		var nodes []string
		for node := range settings {
			nodes = append(nodes, node)
		}

		watching(consulAddr, nodes...)
	}

	return nil
}

func loadByNames(client *consul.Client, names []string) error {
	for i := range names {
		names[i] = path.Join(consul.Redis, names[i])
	}

	return loadConfig(client, names)
}

func loadAll(client *consul.Client) error {
	keys, err := client.GetChildKeys(consul.Redis)
	if err != nil {
		return err
	}

	return loadConfig(client, keys)
}

func loadConfig(client *consul.Client, prefixKeys []string) error {
	load := func(key, instanceName string, isMaster MasterSlave) error {
		val, err := client.Get(key)
		if err != nil {
			return err
		}

		if err = loadNode(client, isMaster, val, instanceName); err != nil {
			return err
		}

		return nil
	}

	for _, key := range prefixKeys {
		instanceName := path.Base(key)
		if err := load(key, instanceName, true); err != nil {
			return err
		}
	}

	return nil
}

func loadNode(client *consul.Client, isMaster MasterSlave, val, instanceName string) error {
	var configs Configs
	if _, err := toml.Decode(val, &configs); err != nil {
		return err
	}

	_redisSettings := make(map[string]*Group, len(settings))
	for k, v := range settings {
		_redisSettings[k] = v
	}

	group, ok := _redisSettings[instanceName]
	if !ok {
		group = &Group{
			Name:       instanceName,
			PoolSize:   configs.PoolSize,
			RedisConns: make([]Conn, 0),
			IsCluster:  configs.IsCluster,
		}

		_redisSettings[instanceName] = group
	} else {
		if group.RefreshSetting {
			group = &Group{
				Name:        instanceName,
				PoolSize:    configs.PoolSize,
				RedisConns:  make([]Conn, 0),
				IsCluster:   configs.IsCluster,
				RefreshPool: true,
			}

			_redisSettings[instanceName] = group
		}
	}

	if len(configs.Master) > 0 {
		for _, master := range configs.Master {
			master.IsMaster = true
			group.RedisConns = append(group.RedisConns, configToConn(master))
		}
	}

	if len(configs.Slave) > 0 {
		for _, slave := range configs.Master {
			slave.IsMaster = false
			group.RedisConns = append(group.RedisConns, configToConn(slave))
		}
	}

	settings = _redisSettings
	return nil
}

func configToConn(conf msConn) Conn {
	return Conn{
		DB:       conf.DB,
		IsMaster: conf.IsMaster,
		ConnStr:  fmt.Sprintf("%s:%v", conf.IP, conf.Port),
	}
}
