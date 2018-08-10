package mongo

import (
	"errors"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/consul"
	"github.com/JREAMLU/j-kit/ext"
	"github.com/hashicorp/consul/api"
	mgo "gopkg.in/mgo.v2"
)

// Config mongo config
type Config struct {
	InstanceName string
	Prefix       string
	URLs         []string
	DBName       string
}

var (
	mgoClients map[string]*mgo.Session
	mutex      sync.Mutex
)

// Watch watch
func Watch(consulAddr string, reloadConfig chan string, names ...string) {
	for i := range names {
		go func(name string) {
			consul.WatchKey(consulAddr, path.Join(consul.MongoDB, name), func(kvPair *api.KVPair) {
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
				mgoClient, err := LoadConfig(consulAddr, false, node)
				if err != nil {
					log.Printf("Failed on mongo LoadConfig, watchedNode: %v, err: %v \r\n", node, err)
					continue
				}

				mutex.Lock()
				mgoClients[node] = mgoClient[node]
				mutex.Unlock()
			}
		}
	}()
}

// Load load mongo
func Load(consulAddr string, isWatching bool, names ...string) error {
	var err error
	mgoClients, err = LoadConfig(consulAddr, isWatching, names...)
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig load config
func LoadConfig(consulAddr string, isWatching bool, names ...string) (map[string]*mgo.Session, error) {
	client, err := consul.NewClient(consul.SetAddress(consulAddr))
	if err != nil {
		return nil, err
	}

	if isWatching {
		watching(consulAddr, names...)
	}

	if len(names) == 0 {
		return loadAll(client)
	}

	return loadByNames(client, names)
}

func loadByNames(client *consul.Client, names []string) (map[string]*mgo.Session, error) {
	for i := range names {
		names[i] = path.Join(consul.MongoDB, names[i])
	}

	return loadConfig(client, names)
}

// GetMongo get mongo session
func GetMongo(instanceName string) *mgo.Session {
	if _, ok := mgoClients[instanceName]; ok {
		return mgoClients[instanceName]
	}

	return nil
}

// GetAllMongo get all mongo session
func GetAllMongo() map[string]*mgo.Session {
	return mgoClients
}

func loadAll(client *consul.Client) (map[string]*mgo.Session, error) {
	keys, err := client.GetChildKeys(consul.MongoDB)
	if err != nil {
		return nil, err
	}

	return loadConfig(client, keys)
}

func loadConfig(client *consul.Client, keys []string) (map[string]*mgo.Session, error) {
	var sessions = make(map[string]*mgo.Session, len(keys))

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

		session, err := registerMongo(config)
		if err != nil {
			return nil, err
		}

		mutex.Lock()
		sessions[instanceName] = session
		mutex.Unlock()
	}

	return sessions, nil
}

func registerMongo(config Config) (*mgo.Session, error) {
	url := ext.StringSplice(config.Prefix, strings.Join(config.URLs, ","))

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("SESSION IS NIL")
	}

	return session.Clone(), nil
}
