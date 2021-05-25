package v2messaging

import (
	"context"
	"fmt"
	kafkaa "leech-service/cqrs/infrastructure/kafkaa"
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

type Consumer struct {
	lock           sync.Mutex
	server         *kafkaa.Server
	messageHandler OnRecieveMessage
}

// Create.
//	Initializes a new instance of the Consumer.
// @param c kafka client.
func New_Consumer(c *kafkaa.Server) *Consumer {
	return &Consumer{
		server: c,
	}
}

func (sr *Consumer) Start(ctx context.Context, onReceiveMessage OnRecieveMessage) {
	sr.lock.Lock()
	sr.messageHandler = onReceiveMessage
	go sr.receiveMessages(ctx)
	sr.lock.Unlock()
}

func (sr *Consumer) receiveMessages(ctx context.Context) {
	for {
		ev := sr.server.Receive(ctx)
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			sr.server.Assign(e)
		case kafka.RevokedPartitions:
			sr.server.Unassign()
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

func (sr *Consumer) Stop(cancelFunc context.CancelFunc) {
	sr.lock.Lock()
	cancelFunc()
	sr.messageHandler = nil
	sr.lock.Unlock()
}

func (sr *Consumer) Complete(ctx context.Context, msg *kafka.Message) error {
	return sr.server.Complete(ctx, msg)
}
