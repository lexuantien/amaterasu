package v2messaging

import (
	"context"
	kafkaa "leech-service/cqrs/infrastructure/kafkaa"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IMessageSender interface {
	Send(ctx context.Context, message *kafka.Message) error // send to message queue
}

type Producer struct {
	client kafkaa.Client
}

// :Create
func New_Producer(c kafkaa.Client) Producer {
	return Producer{
		client: c,
	}
}

func (ts Producer) Send(ctx context.Context, message *kafka.Message) error {
	return ts.client.Send(ctx, message)
}
