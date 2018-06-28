package redis

// Hash redis hash
type Hash struct {
	Structure
}

// NewHash new hash
func NewHash(instanceName, keyPrefixFmt string) Hash {
	return Hash{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Get hash get
func (h *Hash) Get(keySuffix, field string) (string, error) {
	key := h.InitKey(keySuffix)
	return h.String(false, "HGET", key, field)
}
