package handling

import (
	"context"
	"fmt"
	"leech-service/infrastructure/serialization"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Provides basic common processing code for components that handle
// incoming messages from a receiver.
type IMessageProcesser interface {
	Start() // Starts the listener.
	Stop()  // Stops the listener.
}

// Processes incoming commands from the bus and routes them to the appropriate
type CommandProcessor struct {

	// kafka subscriber
	receiver v2messaging.IMessageReceiver
	// serialize kafka message to json - grpc - html ...
	serializer serialization.ISerializer

	// ignore the second request
	started bool

	// uses a peek/lock technique to retrieve a message from a subscription
	lockObject sync.Mutex
	cancelFunc context.CancelFunc

	// dispatch command to specific command handler
	dispatcher *CommandDispatcher
}

//	Create new command processor
// Initializes a new instance of the CommandProcessor
// @param r Kafka receive maeesage
// @param s The serializer to use for the message body.
func New_CommandProcessor(r v2messaging.IMessageReceiver, s serialization.ISerializer) *CommandProcessor {
	return &CommandProcessor{
		receiver:   r,
		serializer: s,
		dispatcher: New_CommandDispatcher(),
	}
}

// Registers the specified command handler.
func (cp *CommandProcessor) Register(commandHandler v2messaging.ICommandHandler, commands ...interface{}) error {
	return cp.dispatcher.Register(commandHandler, commands...)
}

// Processes the message by calling the registered handler.
func (cp *CommandProcessor) processMessage(command interface{}) bool {
	return cp.dispatcher.ProcessMessage(command)
}

func (cp *CommandProcessor) Start() {
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

func (cp *CommandProcessor) Stop() {
	cp.lockObject.Lock()
	if cp.started {
		cp.receiver.Stop(cp.cancelFunc)
		cp.started = false
	}
	cp.lockObject.Unlock()
}

// recieve message from message queue
func (cp *CommandProcessor) onMessageReceived(ctx context.Context, message *kafka.Message) {

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
