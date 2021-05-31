package kafkaa

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Scram struct {
	Al256              bool
	Username, Password string
}

type KafkaConfig struct {
	// use scram to connect
	Scr *Scram

	Brokers string // urls
	Topic   string // topic name

	// ConfigMap is a map containing standard librdkafka configuration properties as documented in:
	// https://github.com/edenhill/librdkafka/tree/master/CONFIGURATION.md
	ConfigMap map[string]interface{}
	configMap *kafka.ConfigMap
}
