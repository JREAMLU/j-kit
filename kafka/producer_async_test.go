package kafka

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var asyncProducer *AsyncProducer

func TestAsyncSend(t *testing.T) {
	newKafkaProducerAsync()

	Convey("Kafka async send test", t, func() {
		Convey("correct", func() {
			err := asyncProducer.Send([]byte("abc"))
			So(err, ShouldBeNil)
		})
	})
}

func newKafkaProducerAsync() {
	asyncProducer = NewAsyncProducer([]string{"127.0.0.1:9092"}, "test_topic")
}
