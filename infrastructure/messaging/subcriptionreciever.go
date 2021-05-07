package messaging

import (
	"context"
	"leech-service/infrastructure/messagebroker"

	"github.com/segmentio/kafka-go"
)

type OnRecieveMessage func(kafka.Message)

type IMessageReceiver interface {
	Start(OnRecieveMessage)
	Stop()
	CommitMessage(context.Context, kafka.Message) error
}

type SubscriptionReciever struct {
	client         *messagebroker.SubscriptionClient
	messageHandler OnRecieveMessage
}

// :Create
func New_SubscriptionReciever(c *messagebroker.SubscriptionClient) *SubscriptionReciever {
	return &SubscriptionReciever{
		client: c,
	}
}

func (sr *SubscriptionReciever) Start(onReceiveMessage OnRecieveMessage) {
	sr.messageHandler = onReceiveMessage
	go func() {
		sr.receiveMessages()
	}()
}

func (sr *SubscriptionReciever) receiveMessages() {
	for {
		msg, err := sr.client.Receive(context.Background())
		if err != nil {
			return
		}
		sr.messageHandler(msg)
	}
}

func (sr *SubscriptionReciever) Stop() {

}

func (sr *SubscriptionReciever) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return sr.client.CommitMessage(ctx, msg)
}
