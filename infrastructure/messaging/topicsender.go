package messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"
)

type IMessageSender interface {
	Send(ctx context.Context, message, clazzType []byte) error // send to message queue
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

func (ts TopicSender) Send(ctx context.Context, message, clazzType []byte) error {
	return ts.client.Send(ctx, message, clazzType)
}
