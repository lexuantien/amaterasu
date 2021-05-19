package v2messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type OnRecieveMessage func(context.Context, kafka.Message)

// Abstracts the behavior of a receiving component that raises
// an event for every received event.
type IMessageReceiver interface {
	// Starts the listener.
	Start(context.Context, OnRecieveMessage)
	// Stops the listener.
	Stop(context.CancelFunc)
	// Complete process message
	Complete(context.Context, kafka.Message) error
}

type SubscriptionReciever struct {
	lock           sync.Mutex
	client         *messagebroker.SubscriptionClient
	messageHandler OnRecieveMessage
}

// Create.
//	Initializes a new instance of the SubscriptionReciever.
// @param c kafka client.
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

		msg, err := sr.client.Receive(ctx)

		if err != nil {

			// Retry the receive loop if there was an error.
			time.Sleep(100 * time.Millisecond)
			continue
		}

		sr.messageHandler(ctx, msg)

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
