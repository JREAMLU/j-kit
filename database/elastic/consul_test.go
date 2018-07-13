package elastic

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddr = "10.200.202.35:8500"
)

func TestLoadConfig(t *testing.T) {
	Convey("load elastic test", t, func() {
		err := Load(consulAddr, false, false, "BGCluster")
		So(err, ShouldBeNil)
		esClient := GetElastic("BGCluster")

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

func TestLoadAll(t *testing.T) {
	Convey("load elastic all test", t, func() {
		err := Load(consulAddr, false, false)
		So(err, ShouldBeNil)
		esClients := GetAllElastic()

		t.Log(err)
		for _, esClient := range esClients {
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
		}
	})
}
