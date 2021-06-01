package kafkaa

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// https://betterprogramming.pub/kafka-acks-explained-c0515b3b707e
// https://www.cloudkarafka.com/blog/apache-kafka-idempotent-producer-avoiding-message-duplication.html
// https://segment.com/blog/exactly-once-delivery/
// https://www.facebook.com/789477804487535/videos/919801674788480

// client is publisher
type Client struct {
	config   KafkaConfig
	producer *kafka.Producer
}

// create new kafka producer configuartion
// config kafka config class, base on confluent-kafka
// 	*Client kafka producer client
// 	error nil if success or fail
/// exclude metadata.broker.list
func New_ProducerConfig(config KafkaConfig) (*Client, error) {

	c := &Client{}
	config.configMap = &kafka.ConfigMap{} // create new config map for confluent-kafka

	// set data
	config.configMap.SetKey("metadata.broker.list", config.Brokers)

	// loop all config then add to
	for k, v := range config.ConfigMap {
		config.configMap.SetKey(k, v)
	}

	// create new producer to publish message
	producer, err := kafka.NewProducer(config.configMap)

	if err != nil {
		return nil, err
	}

	c.config = config
	c.producer = producer

	return c, nil
}

// send message to kafka
// ctx the context
// message kafka message
// 	error nil if success, else
func (c *Client) Send(ctx context.Context, message *kafka.Message) error {

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	// use default topic or custom topic
	if message.TopicPartition.Topic == nil {
		message.TopicPartition.Topic = &c.config.Topic
	}

	// start produce message and waitting
	c.producer.Produce(message, deliveryChan)

	// listen message status response
	e := <-deliveryChan
	m := e.(*kafka.Message)
	close(deliveryChan)

	// TODO log message
	return m.TopicPartition.Error
}
