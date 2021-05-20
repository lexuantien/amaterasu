package messagebroker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type TopicClient struct {
	cl       KafkaConfig
	producer *kafka.Producer
	// writer *kafka.Writer
}

func New_TopicClient(cl KafkaConfig) *TopicClient {
	tc := &TopicClient{}

	tc.cl = cl

	var config *kafka.ConfigMap

	if cl.Scr != nil {
		var al string = "SCRAM-SHA-256"
		if !cl.Scr.Al256 {
			al = "SCRAM-SHA-512"
		}
		config = &kafka.ConfigMap{
			"metadata.broker.list": cl.Brokers,
			"security.protocol":    "SASL_SSL",
			"sasl.mechanisms":      al,
			"sasl.username":        cl.Scr.Username,
			"sasl.password":        cl.Scr.Password,
			"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
			"enable.auto.commit":   false,
			//"debug":                           "generic,broker,security",
		}
	}

	p, err := kafka.NewProducer(config)

	if err != nil {
		panic(err)
	}

	tc.producer = p

	return tc
}

func (c *TopicClient) Send(ctx context.Context, message *kafka.Message) error {

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)
	message.TopicPartition = kafka.TopicPartition{
		Topic:     &c.cl.Topic,
		Partition: kafka.PartitionAny,
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
