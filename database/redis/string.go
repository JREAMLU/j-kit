package redis

import (
	"errors"
	"fmt"

	"github.com/JREAMLU/j-kit/constant"
)

// String redis string
type String struct {
	Structure
}

const (
	_chunkMGETKeys = 10000

	// EXISTS exists
	EXISTS = "EXISTS"
	// GET get
	GET = "GET"
	// SET set
	SET = "SET"
	// SETNX setnx
	SETNX = "SETNX"
	// GETSET getset
	GETSET = "GETSET"
	// SETBIT setbit
	SETBIT = "SETBIT"
	// GETBIT GETBIt
	GETBIT = "GETBIT"
	// SETEX setex
	SETEX = "SETEX"
	// MGET mget
	MGET = "MGET"
	// SETRANGE setrange
	SETRANGE = "SETRANGE"
	// GETRANGE getrange
	GETRANGE = "GETRANGE"
	// STRLEN strlen
	STRLEN = "STRLEN"
	// MSET mset
	MSET = "MSET"
	// MSETNX msetnx
	MSETNX = "MSETNX"
	// PSETEX psetex
	PSETEX = "PSETEX"
	// INCR incr
	INCR = "INCR"
	// INCRBY incrby
	INCRBY = "INCRBY"
	// INCRBYFLOAT incrbyfloat
	INCRBYFLOAT = "INCRBYFLOAT"
	// DECR decr
	DECR = "DECR"
	// DECRBY decrby
	DECRBY = "DECRBY"
	// APPEND append
	APPEND = "APPEND"
)

// NewString new string
func NewString(instanceName, keyPrefixFmt string) String {
	return String{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Exists exists
func (s *String) Exists(keySuffix string) (bool, error) {
	ok, err := s.Int(SLAVE, EXISTS, s.InitKey(keySuffix))
	return ok == 1, err
}

// Get get
func (s *String) Get(keySuffix string) (string, error) {
	return s.String(SLAVE, GET, s.InitKey(keySuffix))
}

// GetInt get int type
func (s *String) GetInt(keySuffix string) (int, error) {
	return s.Int(SLAVE, GET, s.InitKey(keySuffix))
}

// GetInt64 get int64 type
func (s *String) GetInt64(keySuffix string) (int64, error) {
	return s.Int64(SLAVE, GET, s.InitKey(keySuffix))
}

// Set set
func (s *String) Set(keySuffix string, value interface{}, when int) (bool, error) {
	key := s.InitKey(keySuffix)
	if when == constant.NotExists {
		ok, err := s.Int(MASTER, SETNX, key, value)
		if err != nil {
			return false, err
		}

		return ok > 0, nil
	}

	ok, err := s.String(MASTER, SET, key, value)
	if err != nil {
		return false, err
	}

	return ok == OK, nil
}

// GetSet getset
func (s *String) GetSet(keySuffix string) (string, error) {
	return s.String(SLAVE, GETSET, s.InitKey(keySuffix))
}

//SetBit SetBit
func (s *String) SetBit(keySuffix string, offset, value int) (int, error) {
	return s.Int(MASTER, SETBIT, s.InitKey(keySuffix), offset, value)
}

//GetBit GetBit
func (s *String) GetBit(keySuffix string, offset int) (int, error) {
	return s.Int(SLAVE, GETBIT, s.InitKey(keySuffix), offset)
}

// MGet mget
func (s *String) MGet(keySuffix ...string) ([]string, error) {
	keys := make([]interface{}, len(keySuffix))
	for i := range keySuffix {
		keys[i] = s.InitKey(keySuffix[i])
	}

	chunkKeys := sliceChunk(keys, _chunkMGETKeys)
	results := make([]string, 0)
	for _, cKeys := range chunkKeys {
		result, err := s.Strings(SLAVE, MGET, cKeys...)
		if err != nil {
			return nil, err
		}

		results = append(results, result...)
	}

	return results, nil
}

// SetEX setex
func (s *String) SetEX(keySuffix, value string, timespan int) (bool, error) {
	key := s.InitKey(keySuffix)
	ok, err := s.String(MASTER, SETEX, key, timespan, value)
	if err != nil {
		return false, err
	}

	return ok == OK, nil
}

//SetRange setrange
func (s *String) SetRange(keySuffix, value string, offset int) (int, error) {
	return s.Int(MASTER, SETRANGE, s.InitKey(keySuffix), offset, value)
}

//GetRange getrange
func (s *String) GetRange(keySuffix string, start, end int) (string, error) {
	return s.String(SLAVE, GETRANGE, s.InitKey(keySuffix), start, end)
}

//StrLen strlen
func (s *String) StrLen(keySuffix string) (int, error) {
	return s.Int(SLAVE, STRLEN, s.InitKey(keySuffix))
}

// MSet mset
func (s *String) MSet(keySuffix []string, value []interface{}) (bool, error) {
	params, err := s.getParams(keySuffix, value)
	if err != nil {
		return false, err
	}

	reply, err := s.String(MASTER, MSET, params...)
	if err != nil {
		return false, err
	}

	return reply == OK, err
}

// MSetNX msetnx
func (s *String) MSetNX(keySuffix []string, value []interface{}) (int, error) {
	params, err := s.getParams(keySuffix, value)
	if err != nil {
		return constant.ZeroInt, err
	}

	return s.Int(MASTER, MSETNX, params...)
}

func (s *String) getParams(keySuffix []string, value []interface{}) ([]interface{}, error) {
	if len(keySuffix) != len(value) {
		return nil, errors.New("params error: key, value len must be equal")
	}

	params := make([]interface{}, len(keySuffix)*2)
	n := 0
	for i := range keySuffix {
		params[n] = s.InitKey(fmt.Sprint(keySuffix[i]))
		n++
		params[n] = value[i]
		n++
	}

	return params, nil
}

// PSetEX psetex
func (s *String) PSetEX(keySuffix, value string, milliseconds int) (bool, error) {
	ok, err := s.String(MASTER, PSETEX, s.InitKey(keySuffix), milliseconds, value)
	if err != nil {
		return false, err
	}

	return ok == OK, nil
}

// Incr incr
func (s *String) Incr(keySuffix string) (int, error) {
	return s.Int(MASTER, INCR, s.InitKey(keySuffix))
}

// IncrBy incrby
func (s *String) IncrBy(keySuffix string, value int) (int, error) {
	return s.Int(MASTER, INCRBY, s.InitKey(keySuffix), value)
}

// IncrByFloat incrbyfloat
func (s *String) IncrByFloat(keySuffix string, value float64) (float64, error) {
	return s.Float64(MASTER, INCRBYFLOAT, s.InitKey(keySuffix), value)
}

// Decr decr
func (s *String) Decr(keySuffix string) (int, error) {
	return s.Int(MASTER, DECR, s.InitKey(keySuffix))
}

// DecrBy decrby
func (s *String) DecrBy(keySuffix string, value int) (int, error) {
	return s.Int(MASTER, DECRBY, s.InitKey(keySuffix), value)
}

// Append append
func (s *String) Append(keySuffix, value string) (int, error) {
	return s.Int(MASTER, APPEND, s.InitKey(keySuffix), value)
}
