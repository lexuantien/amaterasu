package v2messaging

import (
	"context"
	"fmt"
	kafkaBroker "leech-service/infrastructure/kafka.broker"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OnRecieveMessage func(context.Context, *kafka.Message)

// Abstracts the behavior of a receiving component that raises
// an event for every received event.
type IMessageReceiver interface {
	// Starts the listener.
	Start(context.Context, OnRecieveMessage)
	// Stops the listener.
	Stop(context.CancelFunc)
	// Complete process message
	Complete(context.Context, *kafka.Message) error
}

type SubscriptionReciever struct {
	lock           sync.Mutex
	client         *kafkaBroker.SubscriptionClient
	messageHandler OnRecieveMessage
}

// Create.
//	Initializes a new instance of the SubscriptionReciever.
// @param c kafka client.
func New_SubscriptionReciever(c *kafkaBroker.SubscriptionClient) *SubscriptionReciever {
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
		ev := sr.client.Receive(ctx)
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			sr.client.Assign(e)
		case kafka.RevokedPartitions:
			sr.client.Unassign()
		case *kafka.Message:
			sr.messageHandler(ctx, e)
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			// Retry the receive loop if there was an error
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
}

func (sr *SubscriptionReciever) Stop(cancelFunc context.CancelFunc) {
	sr.lock.Lock()
	cancelFunc()
	sr.messageHandler = nil
	sr.lock.Unlock()
}

func (sr *SubscriptionReciever) Complete(ctx context.Context, msg *kafka.Message) error {
	return sr.client.Complete(ctx, msg)
}
