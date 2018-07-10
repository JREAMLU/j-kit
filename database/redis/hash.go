package redis

// Hash redis hash
type Hash struct {
	Structure
}

const (
	// _maxHMGETFields = 10000
	_chunkHMGETFields = 1
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

	return h.Bool(MASTER, "HDEL", args...)
}

//Exists true: exist false:not exist
func (h *Hash) Exists(keySuffix, field string) (bool, error) {
	key := h.InitKey(keySuffix)
	reply, err := h.Int(MASTER, "HEXISTS", key, field)
	if err != nil {
		return false, err
	}

	return reply != 0, nil
}

// Get hash get
func (h *Hash) Get(keySuffix, field string) (string, error) {
	key := h.InitKey(keySuffix)
	return h.String(SLAVE, "HGET", key, field)
}

// Gets hash gets map
func (h *Hash) Gets(keySuffix string, fields []string) (map[string]string, error) {
	key := h.InitKey(keySuffix)
	result := make(map[string]string)
	chunkFields := sliceChunkString(fields, _chunkHMGETFields)

	for _, cFields := range chunkFields {
		args := append([]interface{}{key}, cFields...)
		reply, err := h.Strings(SLAVE, "HMGET", args...)
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

		subResult, err := h.Ints(false, "HMGET", args...)
		if err != nil {
			return nil, err
		}

		results = append(results, subResult...)
	}

	return results, nil
}
