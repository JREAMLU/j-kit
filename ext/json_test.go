package ext

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJson(t *testing.T) {
	raw := `{"name": "abc",  "age":1}`
	var jsPretty []byte
	Convey("json test", t, func() {
		Convey("pretty", func() {
			var err error
			jsPretty, err = PrettyJSON([]byte(raw))
			So(err, ShouldBeNil)
			So(string(jsPretty), ShouldNotBeNil)
		})

		Convey("minify", func() {
			js, err := Minify(string(jsPretty))
			So(err, ShouldBeNil)
			So(js, ShouldNotBeNil)
		})
	})
}
