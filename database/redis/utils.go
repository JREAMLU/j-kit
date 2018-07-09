package redis

func cutSlice(cut int, src []interface{}) [][]interface{} {
	l := make([][]interface{}, 0)
	start := 0
	offset := cut
	for {
		if start > len(src)-1 {
			break
		}
		//offset += 1
		if offset > len(src) {
			offset = len(src)
		}

		temp := src[start:offset]
		l = append(l, temp)
		start = offset
		offset += cut
	}
	return l
}

func cutStringSlice(cut int, src []string) [][]interface{} {
	l := make([][]interface{}, 0)
	start := 0
	offset := cut
	for {
		if start > len(src)-1 {
			break
		}
		//offset += 1
		if offset > len(src) {
			offset = len(src)
		}

		temp := src[start:offset]
		array := make([]interface{}, len(temp))
		for i, v := range temp {
			array[i] = v
		}
		l = append(l, array)
		start = offset
		offset += cut
	}
	return l
}

func sliceChunkString(sliceStr []string, size int) (chunkslice [][]interface{}) {
	sliceStrLen := len(sliceStr)

	size1 := len(sliceStr) / size
	if size == 0 || sliceStrLen%size > 0 {
		size1++
	}

	chunkSize := (sliceStrLen + size1 - 1) / size1

	slice := make([]interface{}, sliceStrLen)
	for key, val := range sliceStr {
		slice[key] = val
	}

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunkslice = append(chunkslice, slice[i:end])
	}

	return chunkslice
}

func sliceChunkStr(slice []string, size int) (chunkslice [][]string) {
	size1 := len(slice) / size
	if size == 0 || len(slice)%size > 0 {
		size1++
	}

	chunkSize := (len(slice) + size1 - 1) / size1

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunkslice = append(chunkslice, slice[i:end])
	}

	return chunkslice
}
