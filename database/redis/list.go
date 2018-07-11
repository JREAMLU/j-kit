package redis

// List redis list
type List struct {
	Structure
}

const (
	// BLPOP blpop
	BLPOP = "BLPOP"
	// BRPOP brpop
	BRPOP = "BRPOP"
	// BRPOPLPUSH brpoplpush
	BRPOPLPUSH = "BRPOPLPUSH"
	// LINDEX lindex
	LINDEX = "LINDEX"
	// LINSERT linsert
	LINSERT = "LINSERT"
	// BEFORE before
	BEFORE = "BEFORE"
	// AFTER after
	AFTER = "AFTER"
	// LLEN llen
	LLEN = "LLEN"
	// LPOP lpop
	LPOP = "LPOP"
	// LPUSH lpush
	LPUSH = "LPUSH"
	// LPUSHX lpushx
	LPUSHX = "LPUSHX"
	// LRANGE lrange
	LRANGE = "LRANGE"
	// LREM lrem
	LREM = "LREM"
	// LSET lset
	LSET = "LSET"
	// LTRIM ltrim
	LTRIM = "LTRIM"
	// RPOP rpop
	RPOP = "RPOP"
	// RPOPLPUSH rpoplpush
	RPOPLPUSH = "RPOPLPUSH"
	// RPUSH rpush
	RPUSH = "RPUSH"
	// RPUSHX rpushx
	RPUSHX = "RPUSHX"
)

// NewList new list
func NewList(instanceName, keyPrefixFmt string) List {
	return List{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// BLPop blpop
func (l *List) BLPop(keySuffix string, timeout int) ([]string, error) {
	return l.Strings(MASTER, BLPOP, l.InitKey(keySuffix), timeout)
}

// BRPop brpop
func (l *List) BRPop(keySuffix string, timeout int) ([]string, error) {
	return l.Strings(MASTER, BRPOP, l.InitKey(keySuffix), timeout)
}

// BRPopLPush brpoplpush
func (l *List) BRPopLPush(keySuffix, destkey string, timeout int) (string, error) {
	return l.String(MASTER, BRPOPLPUSH, l.InitKey(keySuffix), destkey, timeout)
}

// LIndex lindex
func (l *List) LIndex(keySuffix string, index int) (string, error) {
	return l.String(SLAVE, LINDEX, l.InitKey(keySuffix), index)
}

// InsertBefore insertbefore
func (l *List) InsertBefore(keySuffix string, pivot, value string) (int, error) {
	return l.Int(MASTER, LINSERT, l.InitKey(keySuffix), BEFORE, pivot, value)
}

// InsertAfter insertafter
func (l *List) InsertAfter(keySuffix string, pivot, value string) (int, error) {
	return l.Int(MASTER, LINSERT, l.InitKey(keySuffix), AFTER, pivot, value)
}

// Len len
func (l *List) Len(keySuffix string) (int64, error) {
	return l.Int64(SLAVE, LLEN, l.InitKey(keySuffix))
}

// Pop pop
func (l *List) Pop(keySuffix string) (string, error) {
	return l.String(MASTER, LPOP, l.InitKey(keySuffix))
}

// LPush lpush
func (l *List) LPush(keySuffix, value string) (int64, error) {
	return l.Int64(MASTER, LPUSH, l.InitKey(keySuffix), value)
}

// LPushs lpushs
func (l *List) LPushs(keySuffix string, value ...string) (int64, error) {
	params := make([]interface{}, len(value)+1)
	params[0] = l.InitKey(keySuffix)

	for i := range value {
		params[i+1] = value[i]
	}

	return l.Int64(MASTER, LPUSH, params...)
}

// RPushs rpushs
func (l *List) RPushs(keySuffix string, value ...string) (int64, error) {
	params := make([]interface{}, len(value)+1)
	params[0] = l.InitKey(keySuffix)

	for i := range value {
		params[i+1] = value[i]
	}

	return l.Int64(MASTER, RPUSH, params...)
}

// PushX pushx
func (l *List) PushX(keySuffix, value string) (int64, error) {
	return l.Int64(MASTER, LPUSHX, l.InitKey(keySuffix), value)
}

// LRange lrange
func (l *List) LRange(keySuffix string, start, stop int) ([]string, error) {
	return l.Strings(MASTER, LRANGE, l.InitKey(keySuffix), start, stop)
}

// LRangeInt64 lrangeint64
func (l *List) LRangeInt64(keySuffix string, start, stop int) ([]int64, error) {
	return l.Int64s(MASTER, LRANGE, l.InitKey(keySuffix), start, stop)
}

// Remove remove
func (l *List) Remove(keySuffix, value string, count int) (int64, error) {
	return l.Int64(MASTER, LREM, l.InitKey(keySuffix), count, value)
}

// Set set
func (l *List) Set(keySuffix, value string, index int) (bool, error) {
	reply, err := l.String(MASTER, LSET, l.InitKey(keySuffix), index, value)
	if err != nil {
		return false, err
	}

	return reply == OK, nil
}

// Trim trim
func (l *List) Trim(keySuffix string, start, stop int) (bool, error) {
	reply, err := l.String(MASTER, LTRIM, l.InitKey(keySuffix), start, stop)
	if err != nil {
		return false, err
	}

	return reply == OK, nil
}

// RPop rpop
func (l *List) RPop(keySuffix string) (string, error) {
	return l.String(MASTER, RPOP, l.InitKey(keySuffix))
}

// RPopLPush rpoplpush
func (l *List) RPopLPush(keySuffix, destKey string) (string, error) {
	return l.String(MASTER, RPOPLPUSH, l.InitKey(keySuffix), destKey)
}

// RPushX rpushx
func (l *List) RPushX(keySuffix, value string) (int64, error) {
	return l.Int64(MASTER, RPUSHX, l.InitKey(keySuffix), value)
}
