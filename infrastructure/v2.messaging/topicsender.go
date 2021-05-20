package v2messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IMessageSender interface {
	Send(ctx context.Context, message *kafka.Message) error // send to message queue
}

type TopicSender struct {
	client messagebroker.TopicClient
}

// :Create
func New_TopicSender(c messagebroker.TopicClient) TopicSender {
	return TopicSender{
		client: c,
	}
}

func (ts TopicSender) Send(ctx context.Context, message *kafka.Message) error {
	return ts.client.Send(ctx, message)
}
