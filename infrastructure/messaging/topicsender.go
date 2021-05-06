package messaging

import (
	"context"
	"leech-service/infrastructure/msgqueue"
	"time"
)

type IMessageSender interface {
	Send(ctx context.Context, message []byte, time2live *time.Time) // send to message queue
}

type TopicSender struct {
	client msgqueue.Client
}

func Create(c msgqueue.Client) TopicSender {
	return TopicSender{
		client: c,
	}
}

func (ts TopicSender) Send(ctx context.Context, message []byte, time2live *time.Time) {
	ts.client.Send(ctx, message, time2live)
}
