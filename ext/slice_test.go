package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSliceChunkString(t *testing.T) {
	Convey("Slice Chunk String test", t, func() {
		Convey("odd", func() {
			s := []string{"a", "b", "c", "d"}
			cs := SliceChunkString(s, 2)
			So(len(cs), ShouldEqual, 2)
		})

		Convey("even", func() {
			s := []string{"a", "b", "c", "d", "e"}
			cs := SliceChunkString(s, 2)
			So(len(cs), ShouldEqual, 3)
		})
	})
}

func TestSliceDiff(t *testing.T) {
	Convey("Slice Diff test", t, func() {
		Convey("int64", func() {
			s1 := []int64{1, 2, 3, 4}
			s2 := []int64{1, 2}
			ds := SliceDiffInt64(s1, s2)
			So(len(ds), ShouldEqual, 2)
		})

		Convey("string", func() {
			s1 := []string{"a", "b", "c", "d"}
			s2 := []string{"a", "b"}
			ds := SliceDiffString(s1, s2)
			So(len(ds), ShouldEqual, 2)
		})
	})
}
