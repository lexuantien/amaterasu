package messagebroker

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type TopicClient struct {
	topic  string
	writer *kafka.Writer
}

func New_TopicClient(scr Scram, brokers []string, topic string) *TopicClient {
	client := &TopicClient{}
	client.topic = topic

	var al scram.Algorithm = scram.SHA512
	if scr.Al256 {
		al = scram.SHA256
	}
	mechanism, err := scram.Mechanism(al, scr.Username, scr.Password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
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

func (c *TopicClient) Send(ctx context.Context, message []byte, time2live *time.Time) (bool, error) {
	var err error
	if time2live != nil {
		err = c.writer.WriteMessages(ctx, kafka.Message{
			Topic: c.topic,
			Value: message,
			Time:  *time2live,
		})
	} else {
		err = c.writer.WriteMessages(ctx, kafka.Message{
			Topic: c.topic,
			Value: message,
		})
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
