package kafka

import (
	"strings"
	"testing"

	"github.com/Shopify/sarama"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConsumerKafka(t *testing.T) {
	Convey("Kafka consumer test", t, func() {
		Convey("correct", func() {
			StartConsumeFromKafka("/kafka", "test_group",
				[]string{"test_topic"}, strings.Split("192.168.1.1:2181,192.168.1.2:2181,192.168.1.3:2181", ","), handle)
		})
	})
}

func handle(msg *sarama.ConsumerMessage) error {
	return nil
}
