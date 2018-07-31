package consul

import (
	"fmt"
	"log"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	// KeyNotExist key not exist
	KeyNotExist = "Key(%s) does not exist"
	// DirNotExist dir not exist
	DirNotExist = "Dir(%s) does not exist"
)

// Client consul client
type Client struct {
	config       *api.Config
	consulClient *api.Client
	register     *api.AgentServiceRegistration
}

// NewClient new client
func NewClient(opts ...ClientOptionFunc) (*Client, error) {
	client := &Client{
		config: api.DefaultConfig(),
	}

	for _, opt := range opts {
		opt(client)
	}

	log.Printf("Config Consul Addrs: %v\n", client.config.Address)

	consulClient, err := api.NewClient(client.config)
	if err != nil {
		return nil, err
	}
	client.consulClient = consulClient

	// consulClient.Agent().ServiceRegister(client.register)

	return client, nil
}

// ClientOptionFunc client option func
type ClientOptionFunc func(*Client) error

// SetRegister set register
func SetRegister(reg *api.AgentServiceRegistration) ClientOptionFunc {
	return func(client *Client) error {
		client.register = reg
		return nil
	}
}

// SetAddress set address
func SetAddress(address string) ClientOptionFunc {
	return func(client *Client) error {
		if address != "" {
			client.config.Address = address
		}

		return nil
	}
}

// SetScheme set scheme
func SetScheme(scheme string) ClientOptionFunc {
	return func(client *Client) error {
		if scheme != "" {
			client.config.Scheme = scheme
		}

		return nil
	}
}

// SetDatacenter set dc
func SetDatacenter(datacenter string) ClientOptionFunc {
	return func(client *Client) error {
		if datacenter != "" {
			client.config.Datacenter = datacenter
		}

		return nil
	}
}

// SetHTTPBasicAuth set auth
func SetHTTPBasicAuth(userName, password string) ClientOptionFunc {
	return func(client *Client) error {
		if userName != "" {
			client.config.HttpAuth = &api.HttpBasicAuth{
				Username: userName,
				Password: password,
			}
		}

		return nil
	}
}

// SetWaitTime set wt
func SetWaitTime(waitTime time.Duration) ClientOptionFunc {
	return func(client *Client) error {
		if waitTime > 0 {
			client.config.WaitTime = waitTime
		}
		return nil
	}
}

// SetToken set token
func SetToken(token string) ClientOptionFunc {
	return func(client *Client) error {
		if token != "" {
			client.config.Token = token
		}

		return nil
	}
}

// KV kv client
func (client *Client) KV() *api.KV {
	return client.consulClient.KV()
}

// Deregister deregister
func (client *Client) Deregister(serviceName string) error {
	return client.consulClient.Agent().ServiceDeregister(serviceName)
}

// Put put kv
func (client *Client) Put(key, value string) error {
	pair := &api.KVPair{
		Key:   key,
		Value: []byte(value),
	}

	_, err := client.KV().Put(pair, nil)
	return err
}

// Get get kv
func (client *Client) Get(key string) (string, error) {
	kvPair, _, err := client.KV().Get(key, nil)
	if err != nil {
		return "", err
	}

	if kvPair == nil {
		return "", fmt.Errorf(KeyNotExist, key)
	}

	return string(kvPair.Value), nil
}

// Delete delete kv
func (client *Client) Delete(key string) error {
	_, err := client.KV().Delete(key, nil)

	return err
}

// GetOrDefault get kv, if not value, return default
func (client Client) GetOrDefault(key, defaultValue string) string {
	value, err := client.Get(key)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetInt get int
func (client *Client) GetInt(key string) (int, error) {
	value, err := client.Get(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

// GetInt64 get int64
func (client Client) GetInt64(key string) (int64, error) {
	value, err := client.Get(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(value, 10, 64)
}

// GetOrDefaultInt64 get int64, if not value, return default
func (client Client) GetOrDefaultInt64(key string, defaultValue int64) int64 {
	value, err := client.GetInt64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetFloat64 get float64
func (client *Client) GetFloat64(key string) (float64, error) {
	value, err := client.Get(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}

// GetHostPort get host port
func (client *Client) GetHostPort(key string) (string, string, error) {
	value, err := client.Get(key)
	if err != nil {
		return "", "", err
	}

	return net.SplitHostPort(value)
}

//GetChildKeys get child all keys
func (client *Client) GetChildKeys(keyPrefix string) ([]string, error) {
	if !strings.HasSuffix(keyPrefix, "/") {
		keyPrefix += "/"
	}

	keys, _, err := client.KV().Keys(keyPrefix, "/", nil)
	if err != nil {
		return nil, err
	}

	keyPrefixIndex := -1
	for i := range keys {
		if keys[i] == keyPrefix {
			keyPrefixIndex = i
			break
		}
	}

	if keyPrefixIndex == -1 {
		return keys, nil
	}

	if keyPrefixIndex+1 == len(keys) {
		return keys[:keyPrefixIndex], nil
	}

	if keyPrefixIndex == 0 {
		return keys[1:], nil
	}

	return append(keys[:keyPrefixIndex], keys[keyPrefixIndex+1:]...), nil
}

//GetChildValues get child all keys' value
func (client *Client) GetChildValues(keyPrefix string) (api.KVPairs, error) {
	keys, err := client.GetChildKeys(keyPrefix)
	if err != nil {
		return nil, err
	}

	return client.GetValues(keys)
}

// GetArray get array
// eg.
//          key                     value
// conn/mongodb/goimhistory/1  172.16.9.221:27017
// conn/mongodb/goimhistory/2  172.16.9.222:27017
// keyPrefix=conn/mongodb/goimhistory/
// return ["172.16.9.221:27017","172.16.9.222:27017"]
func (client *Client) GetArray(keyPrefix string) ([]string, error) {
	kvPairs, err := client.GetChildValues(keyPrefix)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(kvPairs))
	for i := range kvPairs {
		result[i] = string(kvPairs[i].Value)
	}

	return result, nil
}

// GetMap get map
// eg.
//             key                         value
// conn/redis/config/master/1/db           0
// conn/redis/config/master/1/ip           172.16.9.221
// conn/redis/config/master/1/poolsize     1
// conn/redis/config/master/1/port         6379
// keyPrefix=conn/redis/config/master/1/
// return map[db:0 ip:172.16.9.221 poolsize:1 port:6379]
func (client *Client) GetMap(keyPrefix string) (map[string]string, error) {
	kvPairs, err := client.GetChildValues(keyPrefix)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for i := range kvPairs {
		result[path.Base(kvPairs[i].Key)] = string(kvPairs[i].Value)
	}

	return result, nil
}

// FunctionIDsTmpl tmpl
const FunctionIDsTmpl = "service/go/%s/functionids"

// GetValues get values
func (client *Client) GetValues(keys []string) (api.KVPairs, error) {
	txn := make(api.KVTxnOps, len(keys))
	for i := range keys {
		txn[i] = &api.KVTxnOp{
			Verb: api.KVGet,
			Key:  keys[i],
		}
	}
	success, resp, _, err := client.KV().Txn(txn, nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Errors) != 0 {
		return nil, fmt.Errorf("%v", resp)
	}

	if !success {
		//TODO
		return nil, nil
	}

	return resp.Results, nil
}

// GetURIAndFunctionIDs get functionids id -> uri
func (client *Client) GetURIAndFunctionIDs(serviceName string) (map[string]string, error) {
	keyPrefix := fmt.Sprintf(FunctionIDsTmpl, serviceName)
	keys, _, err := client.KV().Keys(keyPrefix, "", &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	functionIDs := make(map[string]string)
	kvPairs, err := client.GetValues(keys)
	if err != nil {
		return nil, err
	}

	for i := range kvPairs {
		functionIDs[strings.Replace(kvPairs[i].Key, keyPrefix, "", 1)] = string(kvPairs[i].Value)
	}

	return functionIDs, nil
}
