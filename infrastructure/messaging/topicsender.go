package messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"
	"time"
)

type IMessageSender interface {
	Send(ctx context.Context, message []byte, time2live *time.Time) // send to message queue
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

func (ts TopicSender) Send(ctx context.Context, message []byte, time2live *time.Time) {
	ts.client.Send(ctx, message, time2live)
}
