package messagebroker

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type TopicClient struct {
	topic  string
	writer *kafka.Writer
}

func New_TopicClient(scr Scram, brokers []string, topic string) *TopicClient {
	client := &TopicClient{
		topic: topic,
	}

	var al scram.Algorithm = scram.SHA512
	if scr.Al256 {
		al = scram.SHA256
	}
	mechanism, err := scram.Mechanism(al, scr.Username, scr.Password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		SASLMechanism: mechanism,
	}

	client.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Balancer:     &kafka.Hash{},
		Dialer:       dialer,
		RequiredAcks: int(kafka.RequireOne),
	})

	return client
}

func (c *TopicClient) SetTopic(topic string) {
	c.topic = topic
}

func (c *TopicClient) Send(ctx context.Context, message, clazzType []byte) error {

	return c.writer.WriteMessages(ctx, kafka.Message{
		Topic: c.topic,
		Key:   clazzType,
		Value: message,
	})

}
