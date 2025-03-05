package skafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	DefaultId = "default"

	RoleConsumer = "consumer"
	RoleProducer = "producer"
)

type Config struct {
	// common
	Id               string // optional
	Role             string // optional, default is consumer
	Auto             bool   // optional
	BootstrapServers string // must
	SecurityProtocol string // optional
	SASLUsername     string // optional
	SASLPassword     string // optional
	SASLMechanism    string // optional
	SSLCaLocation    string // optional

	// producer
	Topic                 string                     // must
	RequestRequiredAcks   int                        // optional
	DeliveredCallback     func(kafka.TopicPartition) `json:"-"` // optional
	DeliverFailedCallback func(kafka.TopicPartition) `json:"-"` // optional

	// consumer
	Topics                 []string // must
	GroupId                string   // must
	AutoOffsetReset        string   // optional
	MaxPollIntervalMS      int      // optional
	SessionTimeoutMS       int      // optional
	HeartbeatIntervalMS    int      // optional
	FetchMaxBytes          int      // optional
	MaxPartitionFetchBytes int      // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Role == "" {
		c.Role = RoleConsumer
	}
	if c.AutoOffsetReset == "" {
		c.AutoOffsetReset = "latest"
	}
	if c.MaxPollIntervalMS == 0 {
		c.MaxPollIntervalMS = 10000 // 10s
	}
	if c.SessionTimeoutMS == 0 {
		c.SessionTimeoutMS = 10000 // 10s
	}
	if c.HeartbeatIntervalMS == 0 {
		c.HeartbeatIntervalMS = 3000 // 3s
	}
	if c.FetchMaxBytes == 0 {
		c.FetchMaxBytes = 1024000 // 1M
	}
	if c.MaxPartitionFetchBytes == 0 {
		c.MaxPartitionFetchBytes = 512000 // 500K
	}
	return c
}

func (c Config) Map() *kafka.ConfigMap {
	var m kafka.ConfigMap
	if c.Role == RoleProducer {
		m = kafka.ConfigMap{
			"request.required.acks": c.RequestRequiredAcks,
		}
	} else {
		m = kafka.ConfigMap{
			"group.id":                  c.GroupId,
			"auto.offset.reset":         c.AutoOffsetReset,
			"max.poll.interval.ms":      c.MaxPollIntervalMS,
			"session.timeout.ms":        c.SessionTimeoutMS,
			"heartbeat.interval.ms":     c.HeartbeatIntervalMS,
			"fetch.max.bytes":           c.FetchMaxBytes,
			"max.partition.fetch.bytes": c.MaxPartitionFetchBytes,
		}
	}

	// common config
	m["bootstrap.servers"] = c.BootstrapServers
	if c.SecurityProtocol != "" {
		m["security.protocol"] = c.SecurityProtocol
		switch c.SecurityProtocol {
		case "plaintext":
			// config nothing
		case "sasl_plaintext":
			m["sasl.username"] = c.SASLUsername
			m["sasl.password"] = c.SASLPassword
			m["sasl.mechanism"] = c.SASLMechanism
		case "sasl_ssl":
			m["sasl.username"] = c.SASLUsername
			m["sasl.password"] = c.SASLPassword
			m["sasl.mechanism"] = c.SASLMechanism
			m["ssl.ca.location"] = c.SSLCaLocation
		}
	}

	return &m
}
