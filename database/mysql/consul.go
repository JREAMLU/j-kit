package mysql

import (
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/JREAMLU/j-kit/consul"
	"github.com/JREAMLU/j-kit/ext"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	// Driver mysql driver
	Driver = "mysql"
	// Conn conn
	Conn = "%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s&parseTime=True"
	// Local Sets the location for time
	Local = "Asia%2FShanghai"
	// Charset charset
	Charset = "utf8"
	// READONLY readonly
	READONLY = "ReadOnly"
	// MaxOpenConns max open conns
	MaxOpenConns = 200
	// MaxIdleConns max idle conns
	MaxIdleConns = 60
)

var mutex sync.Mutex

// Config mysql config in consul
type Config struct {
	InstanceName string
	DBName       string
	ReadWrite    readwrite
	ReadOnly     readonly
}

type readwrite struct {
	Server   string
	Password string
	Port     string
	UserID   string
	CharSet  string
}

type readonly struct {
	Server   string
	Password string
	Port     string
	UserID   string
	CharSet  string
}

var gx map[string]*gorm.DB

// Watch watch config
func Watch(consulAddr string, reloadConfig chan string, names ...string) {
	for i := range names {
		go func(name string) {
			consul.WatchKey(consulAddr, path.Join(consul.MYSQL, name), func(kvPair *api.KVPair) {
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
				ngx, err := LoadConfig(consulAddr, false, node)
				if err != nil {
					log.Printf("Failed on mysql LoadConfig, watchedNode: %v, err: %v \r\n", node, err)
					continue
				}

				mutex.Lock()
				gx[node] = ngx[node]
				mutex.Unlock()
			}
		}
	}()
}

// Load load mysql
func Load(consulAddr string, isWatching bool, names ...string) error {
	var err error
	gx, err = LoadConfig(consulAddr, isWatching, names...)
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig load config
func LoadConfig(consulAddr string, isWatching bool, names ...string) (map[string]*gorm.DB, error) {
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

// GetReadOnly get readonly
func GetReadOnly(name string) string {
	return ext.StringSplice(name, READONLY)
}

// GetDB get instance db
func GetDB(name string) *gorm.DB {
	if _, ok := gx[name]; ok {
		return gx[name]
	}

	return nil
}

// GetReadOnlyDB get instance readonly db
func GetReadOnlyDB(name string) *gorm.DB {
	if _, ok := gx[GetReadOnly(name)]; ok {
		return gx[GetReadOnly(name)]
	}

	return nil
}

// GetAllDB get all db
func GetAllDB() map[string]*gorm.DB {
	return gx
}

func loadByNames(client *consul.Client, names []string) (map[string]*gorm.DB, error) {
	for i := range names {
		names[i] = path.Join(consul.MYSQL, names[i])
	}

	return loadConfig(client, names)
}

func loadAll(client *consul.Client) (map[string]*gorm.DB, error) {
	keys, err := client.GetChildKeys(consul.MYSQL)
	if err != nil {
		return nil, err
	}

	return loadConfig(client, keys)
}

func loadConfig(client *consul.Client, keys []string) (map[string]*gorm.DB, error) {
	var dbs = make(map[string]*gorm.DB, len(keys)*2)
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

		// read and write
		rwdb, err := registerDatabase(instanceName, config, true)
		if err != nil {
			return nil, err
		}

		// set pool
		rwdb.DB().SetMaxOpenConns(MaxOpenConns)
		rwdb.DB().SetMaxIdleConns(MaxIdleConns)
		mutex.Lock()
		dbs[instanceName] = rwdb
		mutex.Unlock()

		// readonly
		rdb, err := registerDatabase(instanceName, config, false)
		if err != nil {
			return nil, err
		}

		rdb.DB().SetMaxOpenConns(MaxOpenConns)
		rdb.DB().SetMaxIdleConns(MaxIdleConns)

		mutex.Lock()
		dbs[GetReadOnly(instanceName)] = rdb
		mutex.Unlock()
	}

	return dbs, nil
}

func registerDatabase(name string, config Config, isWrite bool) (*gorm.DB, error) {
	if isWrite {
		var charset string
		if config.ReadWrite.CharSet == "" {
			charset = "utf8"
		} else {
			charset = config.ReadWrite.CharSet
		}

		conn := fmt.Sprintf(Conn, config.ReadWrite.UserID, config.ReadWrite.Password,
			config.ReadWrite.Server, config.ReadWrite.Port, config.DBName, charset, Local)
		db, err := gorm.Open(Driver, conn)

		return db, err
	}

	var charset string
	if config.ReadOnly.CharSet == "" {
		charset = "utf8"
	} else {
		charset = config.ReadOnly.CharSet
	}

	conn := fmt.Sprintf(Conn, config.ReadOnly.UserID, config.ReadOnly.Password,
		config.ReadOnly.Server, config.ReadOnly.Port, config.DBName, charset, Local)
	db, err := gorm.Open(Driver, conn)

	return db, err
}
