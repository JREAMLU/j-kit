package ext

// SliceChunkString slice chunk by string
func SliceChunkString(slice []string, size int) (chunkslice [][]string) {
	size = len(slice) / size
	if size == 0 || len(slice)%size > 0 {
		size = size + 1
	}

	chunkSize := (len(slice) + size - 1) / size

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunkslice = append(chunkslice, slice[i:end])
	}

	return chunkslice
}

// SliceDiffInt64 slice diff by int64
func SliceDiffInt64(slice1, slice2 []int64) (diffslice []int64) {
	for _, v := range slice1 {
		if !inSliceInt64(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}

	return diffslice
}

// SliceDiffString slice diff by string
func SliceDiffString(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if !inSliceString(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}

	return diffslice
}

func inSliceInt64(val int64, slice []int64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func inSliceString(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
