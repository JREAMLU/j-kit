package redis

import (
	"errors"
	"strconv"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/garyburd/redigo/redis"
)

// Hash redis hash
type Hash struct {
	Structure
}

const (
	_chunkHMGETFields = 10000
	_blockSize        = 500

	// HDEL hdel
	HDEL = "HDEL"
	// HEXISTS hexists
	HEXISTS = "HEXISTS"
	// HGET hget
	HGET = "HGET"
	// HMGET hmget
	HMGET = "HMGET"
	// HMSET hmset
	HMSET = "HMSET"
	// OK ok
	OK = "OK"
)

// NewHash new hash
func NewHash(instanceName, keyPrefixFmt string) Hash {
	return Hash{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Delete hash delete
func (h *Hash) Delete(keySuffix string, fields ...interface{}) (bool, error) {
	if len(fields) == 0 {
		return false, nil
	}

	key := h.InitKey(keySuffix)
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	copy(args[1:], fields)

	return h.Bool(MASTER, HDEL, args...)
}

//Exists true: exist false:not exist
func (h *Hash) Exists(keySuffix, field string) (bool, error) {
	key := h.InitKey(keySuffix)
	reply, err := h.Int(MASTER, HEXISTS, key, field)
	if err != nil {
		return false, err
	}

	return reply != 0, nil
}

// Get hash get
func (h *Hash) Get(keySuffix, field string) (string, error) {
	key := h.InitKey(keySuffix)
	return h.String(SLAVE, HGET, key, field)
}

// Gets hash gets map
func (h *Hash) Gets(keySuffix string, fields []string) (map[string]string, error) {
	key := h.InitKey(keySuffix)
	result := make(map[string]string)
	chunkFields := sliceChunkString(fields, _chunkHMGETFields)

	for _, cFields := range chunkFields {
		args := append([]interface{}{key}, cFields...)
		reply, err := h.Strings(SLAVE, HMGET, args...)
		if err != nil {
			return nil, err
		}

		for key := range cFields {
			if reply[key] != "" {
				result[fields[key]] = reply[key]
			}
		}
	}

	return result, nil
}

// GetInts hash get ints
func (h *Hash) GetInts(keySuffix string, fields ...interface{}) ([]int, error) {
	key := h.InitKey(keySuffix)
	chunkFields := sliceChunk(fields, _chunkHMGETFields)
	results := make([]int, 0)
	for _, cFields := range chunkFields {
		args := make([]interface{}, len(cFields)+1)
		args[0] = key
		for i := range cFields {
			args[i+1] = cFields[i]
		}

		reply, err := h.Ints(SLAVE, HMGET, args...)
		if err != nil {
			return nil, err
		}

		results = append(results, reply...)
	}

	return results, nil
}

// GetInt hash get int
func (h *Hash) GetInt(keySuffix, field string) (int, error) {
	key := h.InitKey(keySuffix)
	return h.Int(SLAVE, HGET, key, field)
}

// GetInt64 hash get int64
func (h *Hash) GetInt64(keySuffix, field string) (int64, error) {
	key := h.InitKey(keySuffix)
	return h.Int64(SLAVE, HGET, key, field)
}

// GetInt64s hash get int64s
func (h *Hash) GetInt64s(keySuffix string, fields ...interface{}) ([]int64, error) {
	key := h.InitKey(keySuffix)
	chunkFields := sliceChunk(fields, _chunkHMGETFields)
	results := make([]int64, 0)
	for _, cFields := range chunkFields {
		args := make([]interface{}, len(cFields)+1)
		args[0] = key
		for i := range cFields {
			args[i+1] = cFields[i]
		}

		reply, err := h.Strings(SLAVE, HMGET, args...)
		if err != nil {
			return nil, err
		}

		var result = make([]int64, len(reply))
		for i := range reply {
			if reply[i] == "" {
				continue
			}
			result[i], err = strconv.ParseInt(reply[i], 10, 64)
			if err != nil {
				return nil, err
			}
		}

		results = append(results, result...)
	}

	return results, nil
}

// MSet hash hmset
func (h *Hash) MSet(keySuffix string, fields ...interface{}) (string, error) {
	key := h.InitKey(keySuffix)
	args := append([]interface{}{key}, fields...)

	return h.String(MASTER, HMSET, args...)
}

// MSetStruct hash hmset struct
func (h *Hash) MSetStruct(keySuffix string, p interface{}) (string, error) {
	key := h.InitKey(keySuffix)
	return h.String(MASTER, HMSET, redis.Args{key}.AddFlat(p)...)
}

// MSetSafe hmset batch by safe
// fields muset be even, fields=["field1","value1","field2","value2"]
// blockSize muset be even, eg: blockSize=500, update 250 of 500
func (h *Hash) MSetSafe(keySuffix string, blockSize int, fields ...interface{}) (string, error) {
	if len(fields) == 0 || len(fields)%2 == 1 {
		return constant.EmptyStr, errors.New("error: fields len not zero, fields must be even")
	}

	if blockSize%2 == 1 {
		return constant.EmptyStr, errors.New("error: blockSize must be even")
	}

	key := h.InitKey(keySuffix)
	if blockSize == 0 {
		blockSize = _blockSize
	}

	index := 0
	for {
		if index >= len(fields) {
			return OK, nil
		}

		if blockSize >= len(fields[index:]) {
			args := append([]interface{}{key}, fields[index:]...)
			return h.String(MASTER, HMSET, args...)
		}

		args := append([]interface{}{key}, fields[index:index+blockSize]...)
		_, err := h.String(MASTER, HMSET, args...)
		if err != nil {
			return constant.EmptyStr, err
		}

		index += blockSize
	}
}

// GetAllSafe hash get all by safe, hash scan
func (h *Hash) GetAllSafe(keySuffix string) (map[string]string, error) {
	key := h.InitKey(keySuffix)
	return h.ScanAllMap(key, HSCAN)
}
