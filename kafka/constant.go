package kafka

import "time"

const (
	// OffsetsProcessingTimeout offsets processing timeout second
	OffsetsProcessingTimeout = 10 * time.Second
	// OffsetsCommitInterval offsets commit interval
	OffsetsCommitInterval = 10 * time.Second
	// ZookeeperTimeout zookeeper timeout
	ZookeeperTimeout = 30 * time.Second
	// DialTimeout dial timeout
	DialTimeout = 3 * time.Second
	// ReadTimeout read timeout
	ReadTimeout = 10 * time.Second
	// WriteTimeout write timeout
	WriteTimeout = 10 * time.Second
	// AsyncBuff async producer buff
	AsyncBuff = 1000
)
