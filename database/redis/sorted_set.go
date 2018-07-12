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
	key := s.InitKey(keySuffix)
	return s.Int64(MASTER, "ZINCRBY", key, incr, member)
}

// IncrByFloat64 increment float64
func (s *SortedSet) IncrByFloat64(keySuffix, member string, incr float64) (float64, error) {
	return s.Float64(MASTER, ZINCRBY, s.InitKey(keySuffix), incr, member)
}

// @TODO
//InterStore 给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，
//并将该交集(结果集)储存到 destination 。
//默认情况下，结果集中某个成员的分数值是所有给定集下该成员分数值之和。
func (s *SortedSet) InterStore(destSuffix string, keySuffix ...string) (int64, error) {
	params := make([]interface{}, len(keySuffix)+2)
	params[0] = s.InitKey(destSuffix)
	params[1] = len(keySuffix)
	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}
	return s.Int64(SLAVE, "ZINTERSTORE", params...)
}

// func (s *SortedSet) RangeByRank(keySuffix string, start, stop, order int) ([]string, error) {
// 	key := s.InitKey(keySuffix)
// 	if order == 0 {
// 		return s.Strings(SLAVE, "ZRANGE", key, start, stop)
// 	}
// 	return s.Strings(SLAVE, "ZREVRANGE", key, start, stop)
// }

func (s *SortedSet) RangeWITHSCORES(keySuffix string, start, stop int) (map[string]string, error) {
	key := s.InitKey(keySuffix)
	return s.StringMap(SLAVE, "ZRANGE", key, start, stop, "WITHSCORES")
}

func (s *SortedSet) Range(keySuffix string, start, stop int) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZRANGE", key, start, stop)
}

//RangeByLex 通过字典区间返回有序集合的成员。
func (s *SortedSet) RangeByLex(keySuffix, min, max string) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZRANGEBYLEX", key, min, max)
}

//RangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
func (s *SortedSet) RangeByScore(keySuffix string, min, max interface{}) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZRANGEBYSCORE", key, min, max)
}

func (s *SortedSet) RangeByScoreWITHSCORES(keySuffix string, min, max interface{}) (map[string]string, error) {
	key := s.InitKey(keySuffix)
	return s.StringMap(SLAVE, "ZRANGEBYSCORE", key, min, max, "WITHSCORES")
}

//LexCount 计算有序集合中指定字典区间内成员数量。
func (s *SortedSet) LexCount(keySuffix, min, max string) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(SLAVE, "ZLEXCOUNT", key, min, max)
}

func (s *SortedSet) getParams(key string, members []string, scores []interface{}) ([]interface{}, error) {
	if len(members) != len(scores) {
		return nil, errors.New("params error: members, scores len must be equal")
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

func (s *SortedSet) ScoreInt64(keySuffix, member string) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(SLAVE, "ZSCORE", key, member)
}

func (s *SortedSet) ScoreFloat64(keySuffix, member string) (float64, error) {
	key := s.InitKey(keySuffix)
	return s.Float64(SLAVE, "ZSCORE", key, member)
}

//Rank 返回有序集合中指定成员的索引
func (s *SortedSet) Rank(keySuffix, member string) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(SLAVE, "ZRANK", key, member)
}

//Rem 命令用于移除有序集中的一个或多个成员，不存在的成员将被忽略
func (s *SortedSet) Rem(keySuffix string, members ...string) (int64, error) {
	params := make([]interface{}, len(members)+1)
	params[0] = s.InitKey(keySuffix)
	for i := range members {
		params[i+1] = members[i]
	}
	return s.Int64(MASTER, "ZREM", params...)
}

//RemRangeByLex 移除有序集合中给定的字典区间的所有成员。
func (s *SortedSet) RemRangeByLex(keySuffix string, min, max interface{}) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(MASTER, "ZREMRANGEBYLEX", key, min, max)
}

//RemoveRangeByScore 移除有序集合中给定的排名区间的所有成员
func (s *SortedSet) RemoveRangeByScore(keySuffix string, start, stop int64) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(MASTER, "ZREMRANGEBYSCORE", key, start, stop)
}

//RemRangeByRank 移除有序集中，指定排名(rank)区间内的所有成员
func (s *SortedSet) RemRangeByRank(keySuffix string, start, stop int64) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(MASTER, "ZREMRANGEBYRANK", key, start, stop)
}

//RevRange 返回有序集中指定区间内的成员，通过索引，分数从高到底
func (s *SortedSet) RevRange(keySuffix string, start, stop int64) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZREVRANGE", key, start, stop)
}

//RevRangeWITHSCORES
func (s *SortedSet) RevRangeWITHSCORES(keySuffix string, start, stop int64) (map[string]string, error) {
	key := s.InitKey(keySuffix)
	return s.StringMap(SLAVE, "ZREVRANGE", key, start, stop, "WITHSCORES")
}

func (s *SortedSet) RevRangeWITHSCORES2(keySuffix string, start, stop int64) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZREVRANGE", key, start, stop, "WITHSCORES")
}

//RevRangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
//具有相同分数值的成员按字典序的逆序(reverse lexicographical order )排列。
func (s *SortedSet) RevRangeByScore(keySuffix string, start, stop interface{}) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZREVRANGEBYSCORE", key, start, stop)
}

func (s *SortedSet) RevRangeByScoreWITHSCORES(keySuffix string, start, stop interface{}) (map[string]string, error) {
	key := s.InitKey(keySuffix)
	return s.StringMap(SLAVE, "ZREVRANGEBYSCORE", key, start, stop, "WITHSCORES")
}

func (s *SortedSet) RevRangeByScoreWITHSCORESAndLimitMap(keySuffix string, max, min, offset, count interface{}) (map[string]string, error) {
	key := s.InitKey(keySuffix)
	return s.StringMap(SLAVE, "ZREVRANGEBYSCORE", key, max, min, "WITHSCORES", "LIMIT", offset, count)
}

func (s *SortedSet) RevRangeByScoreWITHSCORESAndLimitSlice(keySuffix string, max, min, offset, count interface{}) ([]string, error) {
	key := s.InitKey(keySuffix)
	return s.Strings(SLAVE, "ZREVRANGEBYSCORE", key, max, min, "WITHSCORES", "LIMIT", offset, count)
}

//RevRank 返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序
func (s *SortedSet) RevRank(keySuffix, member string) (int64, error) {
	key := s.InitKey(keySuffix)
	return s.Int64(SLAVE, "ZREVRANK", key, member)
}

//UnionStore 命令计算给定的一个或多个有序集的并集，
// 其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination
// aggregate:sum|min|max
// aggregate=sum结果集中某个成员的分数值是所有给定集下该成员分数值之和 。
func (s *SortedSet) UnionStore(aggregate, destSuffix string, keySuffix ...string) (int64, error) {
	params := make([]interface{}, len(keySuffix)+4)
	params[0] = s.InitKey(destSuffix)
	params[1] = len(keySuffix)
	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}
	params[len(params)-2] = "AGGREGATE"
	params[len(params)-1] = aggregate
	return s.Int64(MASTER, "ZUNIONSTORE", params...)
}

//UnionStoreByWeights aggregate:sum|min|max
//aggregate=sum结果集中某个成员的分数值是所有给定集下该成员分数值乘以权重值后之和
func (s *SortedSet) UnionStoreByWeights(weights []interface{}, aggregate, destSuffix string, keySuffix ...string) (int64, error) {
	numkeys := len(keySuffix)
	if len(weights) != numkeys || numkeys == 0 {
		return 0, errors.New("weights is empty!")
	}

	params := make([]interface{}, numkeys+len(weights)+5)
	params[0] = s.InitKey(destSuffix)
	params[1] = numkeys
	for i := range keySuffix {
		params[i+2] = s.InitKey(keySuffix[i])
	}
	params[numkeys+2] = "WEIGHTS"
	for i := range weights {
		params[numkeys+i+3] = weights[i]
	}
	params[len(params)-2] = "AGGREGATE"
	params[len(params)-1] = aggregate
	return s.Int64(MASTER, "ZUNIONSTORE", params...)
}

func (s *SortedSet) Scan(keySuffix string, cursor, pageSize int) (int, []string, error) {
	key := s.InitKey(keySuffix)
	return s.Structure.Scan(key, ZSCAN, cursor, pageSize)
}
