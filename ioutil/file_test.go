package ioutil

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	filename = "./testfile.txt"
	dirname  = "./testdir"
)

func init() {
	os.Remove(filename)
	os.RemoveAll(dirname)
}

func TestWriteFile(t *testing.T) {
	Convey("WriteFile test", t, func() {
		Convey("not override, create new", func() {
			err := WriteFile(filename, "abc", false)
			So(err, ShouldBeNil)
		})

		Convey("not override, file exists", func() {
			err := WriteFile(filename, "123", false)
			So(err, ShouldNotBeNil)
		})

		Convey("override", func() {
			err := WriteFile(filename, "abc123", true)
			So(err, ShouldBeNil)
		})
	})
}

func TestReadAll(t *testing.T) {
	Convey("ReadAll test", t, func() {
		Convey("correct", func() {
			content, err := ReadAll(filename)
			So(err, ShouldBeNil)
			So(content, ShouldEqual, "abc123")
		})
	})
}

func TestReadAllBytes(t *testing.T) {
	Convey("ReadAllBytes test", t, func() {
		Convey("correct", func() {
			content, err := ReadAllBytes(filename)
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "abc123")
		})
	})
}

func TestMkdireAll(t *testing.T) {
	Convey("MkdireAll test", t, func() {
		Convey("correct", func() {
			err := MkdireAll(dirname)
			So(err, ShouldBeNil)
		})
	})
}
