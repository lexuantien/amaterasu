package handling

import (
	"context"
	"fmt"
	"leech-service/infrastructure/serialization"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventProcessor struct {

	//
	receiver   v2messaging.IMessageReceiver
	serializer serialization.ISerializer

	//
	started bool

	//
	lockObject sync.Mutex
	cancelFunc context.CancelFunc

	//
	dispatcher *EventDispatcher
}

func New_EventProcessor(r v2messaging.IMessageReceiver, s serialization.ISerializer) *EventProcessor {
	return &EventProcessor{
		receiver:   r,
		serializer: s,
		dispatcher: New_EventDispatcher(),
	}
}

// Registers the specified event handler.
func (cp *EventProcessor) Register(eventHandler v2messaging.IEventHandler, events ...interface{}) error {
	return cp.dispatcher.Register(eventHandler, events...)
}

// Processes the message by calling the registered handler.
func (cp *EventProcessor) processMessage(event interface{}) bool {
	if event != nil {
		return cp.dispatcher.DispatchMessage(event)
	}
	return false
}

func (cp *EventProcessor) Start() {
	// lock
	cp.lockObject.Lock()
	if !cp.started {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cp.cancelFunc = cancelFunc
		cp.receiver.Start(ctx, cp.onMessageReceived)
		cp.started = true
	}
	cp.lockObject.Unlock()
}

func (cp *EventProcessor) Stop() {
	cp.lockObject.Lock()
	if cp.started {
		cp.receiver.Stop(cp.cancelFunc)
		cp.started = false
	}
	cp.lockObject.Unlock()
}

// recieve message from message queue
func (cp *EventProcessor) onMessageReceived(ctx context.Context, message *kafka.Message) {

	cmdClassType, _ := cp.serializer.Deserialize(message.Headers[0].Value, "")

	msg, _ := cp.dispatcher.GetCommandType(cmdClassType.(string))

	msg, _ = cp.serializer.Deserialize(message.Value, msg)

	ok := cp.processMessage(msg)

	// commit msg
	if ok {
		loop := 0
		for {
			if err := cp.receiver.Complete(ctx, message); err != nil {
				fmt.Println(err)
				if loop == 3 {
					// TODO use CB to handle error
					// TODO add to dead letter to commit again
					return
				}
				loop++
			} else {
				break
			}
		}
	}
}
