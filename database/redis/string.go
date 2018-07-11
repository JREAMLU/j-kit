package redis

import "github.com/JREAMLU/j-kit/constant"

// String redis string
type String struct {
	Structure
}

const (
	// EXISTS exists
	EXISTS = "EXISTS"
	// GET get
	GET = "GET"
	// SET set
	SET = "SET"
	// SETNX setnx
	SETNX = "SETNX"
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
