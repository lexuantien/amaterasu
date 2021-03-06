package messaging

import (
	"amaterasu/kafkaa"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IMessageProducer interface {
	Send(ctx context.Context, message *kafka.Message) error // send to message queue
}

type Producer struct {
	client kafkaa.Client
}

//
func New_Producer(c kafkaa.Client) Producer {
	return Producer{
		client: c,
	}
}

func (ts Producer) Send(ctx context.Context, message *kafka.Message) error {
	return ts.client.Send(ctx, message)
}
