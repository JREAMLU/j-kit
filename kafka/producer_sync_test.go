package kafka

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var syncProducer *SyncProducer

func TestSyncSend(t *testing.T) {
	newKafkaProducerSync()

	Convey("Kafka sync send test", t, func() {
		Convey("correct", func() {
			err := asyncProducer.Send([]byte("abc"))
			So(err, ShouldBeNil)
		})
	})
}

func newKafkaProducerSync() {
	var err error
	syncProducer, err = NewSyncProducer("test_topic", []string{"127.0.0.1:9092"}, true)
	if err != nil {
		panic(err)
	}
}
