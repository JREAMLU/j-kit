package ext

// IAbs iabs int
func IAbs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// IAbs32 iabs int32
func IAbs32(x int32) int32 {
	if x < 0 {
		return -x
	}

	return x
}

// IAbs64 iabs int64
func IAbs64(x int64) int64 {
	if x < 0 {
		return -x
	}

	return x
}
