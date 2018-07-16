package elastic

import (
	"log"
	"path"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/consul"
	"github.com/hashicorp/consul/api"
)

// Config elastic in consul
type Config struct {
	InstanceName string
	Debug        bool
	URLs         []string
}

var esCients map[string]*Elastic
var mutex sync.Mutex

// Watch watch config
func Watch(consulAddr string, reloadConfig chan string, names ...string) {
	for i := range names {
		go func(name string) {
			consul.WatchKey(consulAddr, path.Join(consul.ElasticSearch, name), func(kvPair *api.KVPair) {
				reloadConfig <- name
			})
		}(names[i])
	}
}

func watching(consulAddr string, debug bool, names ...string) {
	watchdNode := make(chan string)
	Watch(consulAddr, watchdNode, names...)
	go func() {
		for {
			select {
			case node := <-watchdNode:
				log.Printf("changed: %v \r\n", node)
				nEsclient, err := LoadConfig(consulAddr, false, debug, node)
				if err != nil {
					log.Printf("Failed on mysql LoadConfig, watchedNode: %v, err: %v \r\n", node, err)
					continue
				}

				// @TODO Q
				esCients = nEsclient
			}
		}
	}()
}

// Load load elastic
func Load(consulAddr string, isWatching, debug bool, names ...string) error {
	var err error
	esCients, err = LoadConfig(consulAddr, isWatching, debug, names...)
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig load config
func LoadConfig(consulAddr string, isWatching, debug bool, names ...string) (map[string]*Elastic, error) {
	client, err := consul.NewClient(consul.SetAddress(consulAddr))
	if err != nil {
		return nil, err
	}

	if isWatching {
		watching(consulAddr, debug, names...)
	}

	if len(names) == 0 {
		return loadAll(client)
	}

	return loadByNames(client, names)
}

// GetElastic get elastic
func GetElastic(instanceName string) *Elastic {
	if _, ok := esCients[instanceName]; ok {
		return esCients[instanceName]
	}

	return nil
}

// GetAllElastic get all elastic
func GetAllElastic() map[string]*Elastic {
	return esCients
}

func loadByNames(client *consul.Client, names []string) (map[string]*Elastic, error) {
	for i := range names {
		names[i] = path.Join(consul.ElasticSearch, names[i])
	}

	return loadConfig(client, names)
}

func loadAll(client *consul.Client) (map[string]*Elastic, error) {
	keys, err := client.GetChildKeys(consul.ElasticSearch)
	if err != nil {
		return nil, err
	}

	return loadConfig(client, keys)
}

func loadConfig(client *consul.Client, keys []string) (map[string]*Elastic, error) {
	var ess = make(map[string]*Elastic, len(keys))

	for _, key := range keys {
		instanceName := path.Base(key)
		buf, err := client.Get(key)
		if err != nil {
			return nil, err
		}

		var config Config
		if _, err = toml.Decode(buf, &config); err != nil {
			return nil, err
		}

		es, err := registerElastic(config)
		if err != nil {
			return nil, err
		}

		mutex.Lock()
		ess[instanceName] = es
		mutex.Unlock()
	}

	return ess, nil
}

func registerElastic(config Config) (*Elastic, error) {
	return NewElastic(config.Debug, config.URLs)
}
