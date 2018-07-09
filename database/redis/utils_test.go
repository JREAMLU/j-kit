package redis

import (
	"strconv"
	"testing"
)

func BenchmarkCutStringSlice(b *testing.B) {
	var s []string
	for i := 1; i <= 10000; i++ {
		s = append(s, strconv.Itoa(i))
	}
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cutStringSlice(10, s)
	}
	b.StopTimer()
}

func BenchmarkSliceChunkString(b *testing.B) {
	var s []string
	for i := 1; i <= 10000; i++ {
		s = append(s, strconv.Itoa(i))
	}
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sliceChunkString(s, 10)
	}
	b.StopTimer()
}
