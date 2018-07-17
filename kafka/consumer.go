package kafka

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
)

// StartConsumeFromKafka consumer kafka
func StartConsumeFromKafka(zkroot, groupName string, topics, zookeeper []string, f func(*sarama.ConsumerMessage) error) {
	sarama.Logger = log.New(os.Stdout, "[KAFKA]", log.LstdFlags|log.Lshortfile)
	cg, err := JoinConsumerGroup(zkroot, groupName, topics, zookeeper)
	if err != nil {
		log.Fatal(err)
	}
	defer cg.Close()
	go func() {
		for err := range cg.Errors() {
			log.Println("consumer error", err)
		}
	}()

	for msg := range cg.Messages() {
		if err := f(msg); err != nil {
			log.Println(err, string(msg.Value))
		}
		cg.CommitUpto(msg)
	}
}

// JoinConsumerGroup join conumser
func JoinConsumerGroup(zkroot, groupName string, topics, zookeeper []string) (*consumergroup.ConsumerGroup, error) {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = OffsetsProcessingTimeout
	config.Offsets.CommitInterval = OffsetsCommitInterval
	config.Zookeeper.Chroot = zkroot
	config.Zookeeper.Timeout = ZookeeperTimeout
	return consumergroup.JoinConsumerGroup(groupName, topics, zookeeper, config)
}
