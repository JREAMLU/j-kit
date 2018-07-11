package redis

// Set redis set
type Set struct {
	Structure
}

const (
	// SADD sadd
	SADD = "SADD"
	// SCARD scard
	SCARD = "SCARD"
	// SDIFF sdiff
	SDIFF = "SDIFF"
	// SDIFFSTORE sdiffstore
	SDIFFSTORE = "SDIFFSTORE"
	// SINTER sinter
	SINTER = "SINTER"
	// SINTERSTORE sinterstore
	SINTERSTORE = "SINTERSTORE"
	// SMOVE smove
	SMOVE = "SMOVE"
	// SPOP spop
	SPOP = "SPOP"
	// SRANDMEMBER srandmember
	SRANDMEMBER = "SRANDMEMBER"
	// SUNION sunion
	SUNION = "SUNION"
	// SUNIONSTORE sunionstore
	SUNIONSTORE = "SUNIONSTORE"
	// SISMEMBER sismember
	SISMEMBER = "SISMEMBER"
	// SMEMBERS smembers
	SMEMBERS = "SMEMBERS"
	// SREM srem
	SREM = "SREM"
)

// NewSet new set
func NewSet(instanceName, keyPrefixFmt string) Set {
	return Set{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Add add
func (s *Set) Add(keySuffix string, values ...string) (bool, error) {
	key := s.InitKey(keySuffix)
	members := make([]interface{}, len(values)+1)
	members[0] = key

	for i := range values {
		members[i+1] = values[i]
	}

	reply, err := s.Int64(MASTER, SADD, members...)
	if err != nil {
		return false, err
	}

	return reply > 0, nil
}

// SCard scard
func (s *Set) SCard(keySuffix string) (int64, error) {
	return s.Int64(SLAVE, SCARD, s.InitKey(keySuffix))
}

// SDiff sdiff
func (s *Set) SDiff(key ...interface{}) ([]string, error) {
	return s.Strings(SLAVE, SDIFF, key...)
}

// SDiffStore sdiffstore store destination key1 [key2]
func (s *Set) SDiffStore(key ...interface{}) (int, error) {
	return s.Int(MASTER, SDIFFSTORE, key...)
}

// SInter sinter
func (s *Set) SInter(key ...interface{}) ([]string, error) {
	return s.Strings(SLAVE, SINTER, key...)
}

// SInterStore sinterstore store destination [key1] [key2]
func (s *Set) SInterStore(key ...interface{}) (int64, error) {
	return s.Int64(MASTER, SINTERSTORE, key...)
}

// SMove smove source to destination
func (s *Set) SMove(keySuffix, destKey, member string) (int64, error) {
	sourceKey := s.InitKey(keySuffix)
	return s.Int64(MASTER, SMOVE, sourceKey, destKey, member)
}

// SPop spop
func (s *Set) SPop(keySuffix string) (string, error) {
	return s.String(MASTER, SPOP, s.InitKey(keySuffix))
}

// SRandMember srandmember
func (s *Set) SRandMember(keySuffix string, count int) ([]string, error) {
	return s.Strings(SLAVE, SRANDMEMBER, s.InitKey(keySuffix), count)
}

// SUnion sunion
func (s *Set) SUnion(key1, key2 string) ([]string, error) {
	return s.Strings(SLAVE, SUNION, key1, key2)
}

// UnionStore unionstore store destination
func (s *Set) UnionStore(destKey, key1, key2 string) (int64, error) {
	return s.Int64(MASTER, SUNIONSTORE, destKey, key1, key2)
}

// Sismember sismember
func (s *Set) Sismember(keySuffix, member string) (int, error) {
	return s.Int(SLAVE, SISMEMBER, s.InitKey(keySuffix), member)
}

// SMembers smembers
func (s *Set) SMembers(keySuffix string) ([]string, error) {
	return s.Strings(SLAVE, SMEMBERS, s.InitKey(keySuffix))
}

// Remove remove
func (s *Set) Remove(keySuffix, member string) (bool, error) {
	reply, err := s.Int64(MASTER, SREM, s.InitKey(keySuffix), member)
	if err != nil {
		return false, err
	}

	return reply > 0, nil
}

// Scan scan
func (s *Set) Scan(keySuffix string, cursor, pageSize int) (int, []string, error) {
	return s.Structure.Scan(s.InitKey(keySuffix), SSCAN, cursor, pageSize)
}

// GetAllSafe getallsafe
func (s *Set) GetAllSafe(keySuffix string) ([]string, error) {
	return s.ScanAll(s.InitKey(keySuffix), SSCAN)
}
