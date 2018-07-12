package elastic

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	esURL = `http://127.0.0.1:9200`
)

func TestNewElastic(t *testing.T) {
	Convey("new elastic test", t, func() {
		esClient, err := NewElastic(esURL)
		So(err, ShouldBeNil)
		So(esClient.Code, ShouldEqual, 200)

		t.Log(err)
		t.Log(esClient)
		t.Log("return code: ", esClient.Code)
		t.Log("cluster name: ", esClient.Info.ClusterName)
		t.Log("name: ", esClient.Info.Name)
		t.Log("tag line: ", esClient.Info.TagLine)
		t.Log("build hash: ", esClient.Info.Version.BuildHash)
		t.Log("build snapshot: ", esClient.Info.Version.BuildSnapshot)
		t.Log("build timestamp: ", esClient.Info.Version.BuildTimestamp)
		t.Log("lucene version: ", esClient.Info.Version.LuceneVersion)
		t.Log("elasticsearch version: ", esClient.Info.Version.Number)
	})
}
