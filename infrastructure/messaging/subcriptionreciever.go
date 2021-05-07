package messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type RecieverError struct {
	Err error
	Ctx context.Context
}

type OnRecieveMessage func(context.Context, kafka.Message)

type IMessageReceiver interface {
	Start(context.Context, OnRecieveMessage)
	Stop(context.CancelFunc)
	Complete(context.Context, kafka.Message) error
}

type SubscriptionReciever struct {
	lock           sync.Mutex
	client         *messagebroker.SubscriptionClient
	messageHandler OnRecieveMessage
	errCh          chan RecieverError
}

// :Create
func New_SubscriptionReciever(c *messagebroker.SubscriptionClient) *SubscriptionReciever {
	return &SubscriptionReciever{
		client: c,
	}
}

func (sr *SubscriptionReciever) Start(ctx context.Context, onReceiveMessage OnRecieveMessage) {
	sr.lock.Lock()
	sr.messageHandler = onReceiveMessage
	go sr.receiveMessages(ctx)
	sr.lock.Unlock()
}

func (sr *SubscriptionReciever) receiveMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := sr.client.Receive(ctx)
			if err != nil {
				select {
				case sr.errCh <- RecieverError{Err: err, Ctx: ctx}:
				default:
				}

				// Retry the receive loop if there was an error.
				time.Sleep(time.Second)
				continue
			}

			sr.messageHandler(ctx, msg)
		}

	}

}

func (sr *SubscriptionReciever) Stop(cancelFunc context.CancelFunc) {
	sr.lock.Lock()
	cancelFunc()
	sr.messageHandler = nil
	sr.lock.Unlock()
}

func (sr *SubscriptionReciever) Complete(ctx context.Context, msg kafka.Message) error {
	return sr.client.Complete(ctx, msg)
}
