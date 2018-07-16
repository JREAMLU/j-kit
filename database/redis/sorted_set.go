package redis

import (
	"errors"

	"github.com/JREAMLU/j-kit/constant"
)

// SortedSet redis sortedset
type SortedSet struct {
	Structure
}

const (
	// ZADD zadd
	ZADD = "ZADD"
	// ZCARD zcard
	ZCARD = "ZCARD"
	// ZCOUNT zcount
	ZCOUNT = "ZCOUNT"
	// ZINCRBY zincrby
	ZINCRBY = "ZINCRBY"
	// ZINTERSTORE zinterstore
	ZINTERSTORE = "ZINTERSTORE"
	// ZRANGE zrange
	ZRANGE = "ZRANGE"
	// ZRANGEBYLEX zrangebylex
	ZRANGEBYLEX = "ZRANGEBYLEX"
	// ZRANGEBYSCORE zrangebyscore
	ZRANGEBYSCORE = "ZRANGEBYSCORE"
	// ZLEXCOUNT zlexcount
	ZLEXCOUNT = "ZLEXCOUNT"
	// ZSCORE zscore
	ZSCORE = "ZSCORE"
	// ZRANK zrank
	ZRANK = "ZRANK"
	// ZREM zrem
	ZREM = "ZREM"
	// ZREMRANGEBYLEX zremrangebylex
	ZREMRANGEBYLEX = "ZREMRANGEBYLEX"
	// ZREMRANGEBYSCORE zremrangebyscore
	ZREMRANGEBYSCORE = "ZREMRANGEBYSCORE"
	// ZREMRANGEBYRANK zremrangebyrank
	ZREMRANGEBYRANK = "ZREMRANGEBYRANK"
	// ZREVRANGE zrevrange
	ZREVRANGE = "ZREVRANGE"
	// ZREVRANGEBYSCORE zrevrangebyscore
	ZREVRANGEBYSCORE = "ZREVRANGEBYSCORE"
	// ZREVRANK zrevrank
	ZREVRANK = "ZREVRANK"
	// AGGREGATE aggregate
	AGGREGATE = "AGGREGATE"
	// ZUNIONSTORE zunionstore
	ZUNIONSTORE = "ZUNIONSTORE"
)

// NewSortedSet new sortedset
func NewSortedSet(instanceName, keyPrefixFmt string) SortedSet {
	return SortedSet{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Add add
func (s *SortedSet) Add(keySuffix, member string, score interface{}) (int64, error) {
	return s.Int64(MASTER, ZADD, s.InitKey(keySuffix), score, member)
}

// AddNX NX: Don't update already existing elements. Always add new elements.
func (s *SortedSet) AddNX(keySuffix, member string, score interface{}) (int64, error) {
	return s.Int64(MASTER, ZADD, s.InitKey(keySuffix), NX, score, member)
}

// AddMulti add multi
func (s *SortedSet) AddMulti(keySuffix string, members []string, score []interface{}) (int64, error) {
	params, err := s.getParams(s.InitKey(keySuffix), members, score)
	if err != nil {
		return constant.ZeroInt64, err
	}

	return s.Int64(MASTER, ZADD, params...)
}

//Card card
func (s *SortedSet) Card(keySuffix string) (int64, error) {
	return s.Int64(SLAVE, ZCARD, s.InitKey(keySuffix))
}

//Count zcount
func (s *SortedSet) Count(keySuffix string, min, max interface{}) (int64, error) {
	return s.Int64(SLAVE, ZCOUNT, s.InitKey(keySuffix), min, max)
}

// IncrByInt64 increment int64
func (s *SortedSet) IncrByInt64(keySuffix, member string, incr int64) (int64, error) {
	return s.Int64(MASTER, ZINCRBY, s.InitKey(keySuffix), incr, member)
}

// IncrByFloat64 increment float64
func (s *SortedSet) IncrByFloat64(keySuffix, member string, incr float64) (float64, error) {
	return s.Float64(MASTER, ZINCRBY, s.InitKey(keySuffix), incr, member)
}

// InterStore interstore destination numkeys key...
func (s *SortedSet) InterStore(destSuffix string, keySuffix ...string) (int64, error) {
	params := make([]interface{}, len(keySuffix)+2)
	params[0] = s.InitKey(destSuffix)
	params[1] = len(keySuffix)

	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}

	return s.Int64(SLAVE, ZINTERSTORE, params...)
}

// RangeWs zrange withscores
func (s *SortedSet) RangeWs(keySuffix string, start, stop int) (map[string]string, error) {
	return s.StringMap(SLAVE, ZRANGE, s.InitKey(keySuffix), start, stop, WITHSCORES)
}

// Range zrange
func (s *SortedSet) Range(keySuffix string, start, stop int) ([]string, error) {
	return s.Strings(SLAVE, ZRANGE, s.InitKey(keySuffix), start, stop)
}

// RangeByLex rangebylex
func (s *SortedSet) RangeByLex(keySuffix, min, max string) ([]string, error) {
	return s.Strings(SLAVE, ZRANGEBYLEX, s.InitKey(keySuffix), min, max)
}

// RangeByScore zrangebyscore
func (s *SortedSet) RangeByScore(keySuffix string, min, max interface{}) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, ZRANGEBYSCORE, key, min, max)
}

// RangeByScoreWs rangebyscore withscores
func (s *SortedSet) RangeByScoreWs(keySuffix string, min, max interface{}) (map[string]string, error) {
	return s.StringMap(SLAVE, ZRANGEBYSCORE, s.InitKey(keySuffix), min, max, WITHSCORES)
}

// LexCount zlexcount
func (s *SortedSet) LexCount(keySuffix, min, max string) (int64, error) {
	return s.Int64(SLAVE, ZLEXCOUNT, s.InitKey(keySuffix), min, max)
}

// ScoreInt64 zscore int64
func (s *SortedSet) ScoreInt64(keySuffix, member string) (int64, error) {
	return s.Int64(SLAVE, ZSCORE, s.InitKey(keySuffix), member)
}

// ScoreFloat64 zscore float64
func (s *SortedSet) ScoreFloat64(keySuffix, member string) (float64, error) {
	return s.Float64(SLAVE, ZSCORE, s.InitKey(keySuffix), member)
}

// Rank rank
func (s *SortedSet) Rank(keySuffix, member string) (int64, error) {
	return s.Int64(SLAVE, ZRANK, s.InitKey(keySuffix), member)
}

// Remove zrem
func (s *SortedSet) Remove(keySuffix string, members ...string) (int64, error) {
	params := make([]interface{}, len(members)+1)
	params[0] = s.InitKey(keySuffix)

	for i := range members {
		params[i+1] = members[i]
	}

	return s.Int64(MASTER, ZREM, params...)
}

// RemoveRangeByLex remrangebylex
func (s *SortedSet) RemoveRangeByLex(keySuffix string, min, max interface{}) (int64, error) {
	return s.Int64(MASTER, ZREMRANGEBYLEX, s.InitKey(keySuffix), min, max)
}

// RemoveRangeByScore remrangebyscore
func (s *SortedSet) RemoveRangeByScore(keySuffix string, start, stop int64) (int64, error) {
	return s.Int64(MASTER, ZREMRANGEBYSCORE, s.InitKey(keySuffix), start, stop)
}

// RemoveRangeByRank remrangebyrank
func (s *SortedSet) RemoveRangeByRank(keySuffix string, start, stop int64) (int64, error) {
	return s.Int64(MASTER, ZREMRANGEBYRANK, s.InitKey(keySuffix), start, stop)
}

// RevRange revrange
func (s *SortedSet) RevRange(keySuffix string, start, stop int64) ([]string, error) {
	return s.Strings(SLAVE, ZREVRANGE, s.InitKey(keySuffix), start, stop)
}

// RevRangeWs revrange withscores
func (s *SortedSet) RevRangeWs(keySuffix string, start, stop int64) (map[string]string, error) {
	return s.StringMap(SLAVE, ZREVRANGE, s.InitKey(keySuffix), start, stop, WITHSCORES)
}

// RevRangesWs revrange slice withscores
func (s *SortedSet) RevRangesWs(keySuffix string, start, stop int64) ([]string, error) {
	return s.Strings(SLAVE, ZREVRANGE, s.InitKey(keySuffix), start, stop, WITHSCORES)
}

// RevRangeByScore revrangebyscore
func (s *SortedSet) RevRangeByScore(keySuffix string, start, stop interface{}) ([]string, error) {
	return s.Strings(SLAVE, ZREVRANGEBYSCORE, s.InitKey(keySuffix), start, stop)
}

// RevRangeByScoreWs revrangebyscore withscores
func (s *SortedSet) RevRangeByScoreWs(keySuffix string, start, stop interface{}) (map[string]string, error) {
	return s.StringMap(SLAVE, ZREVRANGEBYSCORE, s.InitKey(keySuffix), start, stop, WITHSCORES)
}

// RevRangeByScoreWsPage revrangebyscore withscores pagination
func (s *SortedSet) RevRangeByScoreWsPage(keySuffix string, max, min, offset, count interface{}) (map[string]string, error) {
	return s.StringMap(SLAVE, ZREVRANGEBYSCORE, s.InitKey(keySuffix), max, min, WITHSCORES, LIMIT, offset, count)
}

// RevRangeByScoresWsPage revrangebyscore withscores pagination slice
func (s *SortedSet) RevRangeByScoresWsPage(keySuffix string, max, min, offset, count interface{}) ([]string, error) {
	return s.Strings(SLAVE, ZREVRANGEBYSCORE, s.InitKey(keySuffix), max, min, WITHSCORES, LIMIT, offset, count)
}

// RevRank zrevrank
func (s *SortedSet) RevRank(keySuffix, member string) (int64, error) {
	return s.Int64(SLAVE, ZREVRANK, s.InitKey(keySuffix), member)
}

// UnionStore zunionstore aggregate:sum|min|max
func (s *SortedSet) UnionStore(aggregate, destSuffix string, keySuffix ...string) (int64, error) {
	params := make([]interface{}, len(keySuffix)+4)
	params[0] = s.InitKey(destSuffix)
	params[1] = len(keySuffix)

	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}

	params[len(params)-2] = AGGREGATE
	params[len(params)-1] = aggregate

	return s.Int64(MASTER, ZUNIONSTORE, params...)
}

// UnionStoreByWeights unionstore by weights aggregate:sum|min|max
func (s *SortedSet) UnionStoreByWeights(weights []interface{}, aggregate, destSuffix string, keySuffix ...string) (int64, error) {
	numkeys := len(keySuffix)
	if len(weights) != numkeys || numkeys == 0 {
		return constant.ZeroInt64, errors.New("WEIGHTS OR KEY MUST BE NOT EMPTY")
	}

	params := make([]interface{}, numkeys+len(weights)+5)
	params[0] = s.InitKey(destSuffix)
	params[1] = numkeys

	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}

	params[numkeys+2] = WEIGHTS
	for i := range weights {
		params[numkeys+i+3] = weights[i]
	}

	params[len(params)-2] = AGGREGATE
	params[len(params)-1] = aggregate

	return s.Int64(MASTER, ZUNIONSTORE, params...)
}

// Scan scan
func (s *SortedSet) Scan(keySuffix string, cursor, pageSize int) (int, []string, error) {
	return s.Structure.Scan(s.InitKey(keySuffix), ZSCAN, cursor, pageSize)
}

func (s *SortedSet) getParams(key string, members []string, scores []interface{}) ([]interface{}, error) {
	if len(members) != len(scores) {
		return nil, errors.New("MEMBERS, SCORES LEN MUST BE EQUAL")
	}

	params := make([]interface{}, len(members)*2+1)
	params[0] = key
	n := 1
	for i := range members {
		params[n] = scores[i]
		n++
		params[n] = s.InitKey(members[i])
		n++
	}

	return params, nil
}
