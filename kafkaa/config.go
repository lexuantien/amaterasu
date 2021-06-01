package kafkaa

import "github.com/confluentinc/confluent-kafka-go/kafka"

type KafkaConfig struct {
	// urls server handle kafka
	Brokers string
	// topic name
	Topic string
	// ConfigMap is a map containing standard librdkafka configuration properties as documented in:
	// https://github.com/edenhill/librdkafka/tree/master/CONFIGURATION.md
	ConfigMap map[string]interface{} // for ui
	// for confluent kafka
	configMap *kafka.ConfigMap
}
