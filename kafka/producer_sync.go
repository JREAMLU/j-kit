package kafka

import (
	"errors"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

// SyncProducer producer
type SyncProducer struct {
	syncProducer sarama.SyncProducer
	Topic        string
	debug        bool
	json         jsoniter.API
}

// NewSyncProducer new producer
func NewSyncProducer(topic string, kafkaAddrs []string, debug bool) (*SyncProducer, error) {
	producer := &SyncProducer{
		Topic: topic,
		debug: debug,
		json:  jsoniter.ConfigCompatibleWithStandardLibrary,
	}

	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Net.DialTimeout = DialTimeout
	config.Net.ReadTimeout = ReadTimeout
	config.Net.WriteTimeout = WriteTimeout

	syncProducer, err := sarama.NewSyncProducer(kafkaAddrs, config)
	if err != nil {
		return nil, err
	}

	producer.syncProducer = syncProducer

	return producer, nil
}

// Close close
func (producer *SyncProducer) Close() error {
	if producer.syncProducer == nil {
		return errors.New("PRODUCER IS NIL")
	}

	return producer.syncProducer.Close()
}

// Send send
func (producer *SyncProducer) Send(msg []byte) (int32, int64, error) {
	return producer.syncProducer.SendMessage(&sarama.ProducerMessage{Topic: producer.Topic, Value: sarama.ByteEncoder(msg)})
}

// SendWithTopic Send With Topic
func (producer *SyncProducer) SendWithTopic(topic string, msg []byte) (int32, int64, error) {
	return producer.syncProducer.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msg)})
}

// SendObjectsWithTopic Send Objects With Topic
func (producer *SyncProducer) SendObjectsWithTopic(topic string, msgs []interface{}) error {
	producerMessages := make([]*sarama.ProducerMessage, len(msgs))

	for i := range msgs {
		buf, err := producer.json.Marshal(msgs[i])
		if err != nil {
			return err
		}
		producerMessages[i] = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(buf),
		}
	}

	return producer.syncProducer.SendMessages(producerMessages)
}

// SendMessagesWithTopic Send Messages With Topic
func (producer *SyncProducer) SendMessagesWithTopic(topic string, msgs [][]byte) error {
	producerMessages := make([]*sarama.ProducerMessage, len(msgs))
	for i := range msgs {
		producerMessages[i] = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msgs[i]),
		}
	}

	return producer.syncProducer.SendMessages(producerMessages)
}
