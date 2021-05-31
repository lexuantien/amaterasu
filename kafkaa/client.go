package kafkaa

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// https://betterprogramming.pub/kafka-acks-explained-c0515b3b707e
// https://www.cloudkarafka.com/blog/apache-kafka-idempotent-producer-avoiding-message-duplication.html
// https://segment.com/blog/exactly-once-delivery/
// https://www.facebook.com/789477804487535/videos/919801674788480
type Client struct {
	config   KafkaConfig
	producer *kafka.Producer
}

func New_kafkaa(config KafkaConfig) *Client {
	c := &Client{}
	config.configMap = &kafka.ConfigMap{}

	config.configMap.SetKey("metadata.broker.list", config.Brokers)

	if config.Scr != nil {

		var al string = "SCRAM-SHA-256"
		if !config.Scr.Al256 {
			al = "SCRAM-SHA-512"
		}

		config.configMap.SetKey("security.protocol", "SASL_SSL")
		config.configMap.SetKey("sasl.username", config.Scr.Username)
		config.configMap.SetKey("sasl.password", config.Scr.Password)
		config.configMap.SetKey("sasl.mechanisms", al)

	}

	for k, v := range config.ConfigMap {
		config.configMap.SetKey(k, v)
	}

	p, err := kafka.NewProducer(config.configMap)

	if err != nil {
		panic(err)
	}

	c.config = config
	c.producer = p

	return c
}

func (c *Client) Send(ctx context.Context, message *kafka.Message) error {

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	if message.TopicPartition.Topic == nil {
		message.TopicPartition.Topic = &c.config.Topic
	}

	c.producer.Produce(message, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	close(deliveryChan)

	// TODO log message
	// if m.TopicPartition.Error != nil {
	// 	fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	// } else {
	// 	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
	// 		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	// }

	return m.TopicPartition.Error
}
