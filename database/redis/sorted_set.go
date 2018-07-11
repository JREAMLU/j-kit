package redis

// SortedSet redis sortedset
type SortedSet struct {
	Structure
}

// NewSortedSet new sortedset
func NewSortedSet(instanceName, keyPrefixFmt string) SortedSet {
	return SortedSet{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}
