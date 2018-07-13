package elastic

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	esURL = `http://10.200.202.36:9200/`
)

func TestNewElastic(t *testing.T) {
	Convey("new elastic test", t, func() {
		esClient, err := NewElastic(false, []string{esURL})
		So(err, ShouldBeNil)
		So(esClient.Codes[0], ShouldEqual, 200)

		t.Log(err)
		for i := range esClient.Codes {
			t.Log("return code: ", esClient.Codes[i])
			t.Log("cluster name: ", esClient.Infos[i].ClusterName)
			t.Log("name: ", esClient.Infos[i].Name)
			t.Log("tag line: ", esClient.Infos[i].TagLine)
			t.Log("build hash: ", esClient.Infos[i].Version.BuildHash)
			t.Log("build snapshot: ", esClient.Infos[i].Version.BuildSnapshot)
			t.Log("build timestamp: ", esClient.Infos[i].Version.BuildTimestamp)
			t.Log("lucene version: ", esClient.Infos[i].Version.LuceneVersion)
			t.Log("elasticsearch version: ", esClient.Infos[i].Version.Number)
		}
	})
}
