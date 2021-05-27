package handling

import (
	"amaterasu/cqrs/infrastructure/serialization"
	v2messaging "amaterasu/cqrs/infrastructure/v2.messaging"
	"context"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventProcessor struct {

	//
	consumer   v2messaging.IMessageConsumer
	serializer serialization.ISerializer

	//
	started bool

	//
	lockObject sync.Mutex
	cancelFunc context.CancelFunc

	//
	dispatcher *EventDispatcher
}

func New_EventProcessor(r v2messaging.IMessageConsumer, s serialization.ISerializer) *EventProcessor {
	return &EventProcessor{
		consumer:   r,
		serializer: s,
		dispatcher: New_EventDispatcher(),
	}
}

// Registers the specified event handler.
func (cp *EventProcessor) Register(eventHandler v2messaging.IEventHandler) error {
	return cp.dispatcher.Register(eventHandler)
}

// Processes the message by calling the registered handler.
func (cp *EventProcessor) processMessage(event v2messaging.Envelope) bool {
	return cp.dispatcher.DispatchMessage(event) == nil
}

func (cp *EventProcessor) Start() {
	// lock
	cp.lockObject.Lock()
	if !cp.started {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cp.cancelFunc = cancelFunc
		cp.consumer.Start(ctx, cp.onMessageReceived)
		cp.started = true
	}
	cp.lockObject.Unlock()
}

func (cp *EventProcessor) Stop() {
	cp.lockObject.Lock()
	if cp.started {
		cp.consumer.Stop(cp.cancelFunc)
		cp.started = false
	}
	cp.lockObject.Unlock()
}

// recieve message from message queue
func (cp *EventProcessor) onMessageReceived(ctx context.Context, message *kafka.Message) error {
	return OnMessageReceivedHandler(ctx, message, cp.processMessage, cp.consumer.Complete)
}
