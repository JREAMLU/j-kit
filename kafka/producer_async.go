package kafka

import (
	"errors"
	"log"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

// AsyncProducer producer
type AsyncProducer struct {
	topic         string
	data          chan *sarama.ProducerMessage
	KafkaProducer sarama.AsyncProducer
	json          jsoniter.API
}

// NewAsyncProducer new producer
func NewAsyncProducer(kafkaAddr []string, topic string) *AsyncProducer {
	producer := &AsyncProducer{
		data:  make(chan *sarama.ProducerMessage, AsyncBuff),
		topic: topic,
		json:  jsoniter.ConfigCompatibleWithStandardLibrary,
	}

	if err := producer.initKafka(kafkaAddr); err != nil {
		panic(err)
	}

	go producer.send()

	return producer
}

//Send channel <-
func (producer *AsyncProducer) Send(p []byte) error {
	select {
	case producer.data <- &sarama.ProducerMessage{
		Topic: producer.topic,
		Value: sarama.ByteEncoder(p),
	}:
		return nil
	default:
		return errors.New("KAFKA CHANNEL IS FULL")
	}
}

//SendByWait wait channel
func (producer *AsyncProducer) SendByWait(p []byte) error {
	producer.KafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: producer.topic,
		Value: sarama.ByteEncoder(p),
	}

	return nil
}

//SendObject json and send
func (producer *AsyncProducer) SendObject(obj interface{}) error {
	p, err := producer.json.Marshal(obj)
	if err != nil {
		return err
	}

	return producer.Send(p)
}

// SendObjectWithTopic Send Object With Topic
func (producer *AsyncProducer) SendObjectWithTopic(topic string, obj interface{}) error {
	p, err := producer.json.Marshal(obj)
	if err != nil {
		return err
	}
	return producer.SendWithTopic(topic, p)
}

//SendWithTopic send with topic
func (producer *AsyncProducer) SendWithTopic(topic string, p []byte) error {
	select {
	case producer.data <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(p),
	}:
		return nil
	default:
		return errors.New("KAFKA CHANNEL IS FULL")
	}
}

// SendWithKey send with key
func (producer *AsyncProducer) SendWithKey(key string, p []byte) error {
	select {
	case producer.data <- &sarama.ProducerMessage{
		Topic: producer.topic,
		Key:   sarama.ByteEncoder([]byte(key)),
		Value: sarama.ByteEncoder(p),
	}:
		return nil
	default:
		return errors.New("KAFKA CHANNEL IS FULL")
	}
}

// SendWithKeyAndTopic send with key topic
func (producer *AsyncProducer) SendWithKeyAndTopic(key, topic string, p []byte) error {
	select {
	case producer.data <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder([]byte(key)),
		Value: sarama.ByteEncoder(p),
	}:
		return nil
	default:
		return errors.New("KAFKA CHANNEL IS FULL")
	}
}

func (producer *AsyncProducer) send() {
	var (
		ok  bool
		msg *sarama.ProducerMessage
	)

	for {
		if msg, ok = <-producer.data; !ok {
			break
		}
		producer.KafkaProducer.Input() <- msg
	}
}

func (producer *AsyncProducer) initKafka(kafkaAddrs []string) (err error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true

	config.Net.DialTimeout = DialTimeout
	config.Net.ReadTimeout = ReadTimeout
	config.Net.WriteTimeout = WriteTimeout

	producer.KafkaProducer, err = sarama.NewAsyncProducer(kafkaAddrs, config)

	go producer.handleError()

	return err
}

func (producer *AsyncProducer) handleError() {
	var (
		err *sarama.ProducerError
		ok  bool
	)

	for {
		err, ok = <-producer.KafkaProducer.Errors()
		if err != nil {
			log.Printf("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)\n", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}
		if !ok {
			log.Printf("producer ProducerError has be closed, break the handleError goroutine")
			return
		}
	}
}

// Close close
func (producer *AsyncProducer) Close() error {
	return producer.KafkaProducer.Close()
}
