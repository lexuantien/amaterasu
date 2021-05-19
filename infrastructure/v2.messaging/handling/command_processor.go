package handling

import (
	"context"
	"leech-service/infrastructure/serialization"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Provides basic common processing code for components that handle
// incoming messages from a receiver.
type IMessageProcesser interface {
	Start() // Starts the listener.
	Stop()  // Stops the listener.
}

// Initializes a new instance
type CommandProcessor struct {

	// kafka subscriber
	reciever v2messaging.IMessageReceiver
	// serialize kafka message to json - grpc - html ...
	serializer serialization.ISerializer

	// ignore the second request
	started bool

	// uses a peek/lock technique to retrieve a message from a subscription
	lock       sync.Mutex
	cancelFunc context.CancelFunc

	// dispatch command to specific command handler
	dispatcher *CommandDispatcher
}

//	Create new command processor
func New_CommandProcessor(r v2messaging.IMessageReceiver, s serialization.ISerializer) *CommandProcessor {
	return &CommandProcessor{
		reciever:   r,
		serializer: s,
		dispatcher: New_CommandDispatcher(),
	}
}

func (cp *CommandProcessor) Register(commandHandler v2messaging.ICommandHandler, commands ...interface{}) error {
	return cp.dispatcher.Register(commandHandler, commands...)
}

func (cp *CommandProcessor) processMessage(command interface{}) {
	cp.dispatcher.Dispatch(command)
}

func (cp *CommandProcessor) Start() {
	// lock
	cp.lock.Lock()
	if !cp.started {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cp.cancelFunc = cancelFunc
		cp.reciever.Start(ctx, cp.onMessageReceived)
		cp.started = true
	}
	cp.lock.Unlock()
}

func (cp *CommandProcessor) Stop() {
	cp.lock.Lock()
	if cp.started {
		cp.reciever.Stop(cp.cancelFunc)
		cp.started = false
	}
	cp.lock.Unlock()
}

// recieve message from message queue
func (cp *CommandProcessor) onMessageReceived(ctx context.Context, message kafka.Message) {

	cmdClassType, _ := cp.serializer.Deserialize(message.Key, "")

	msg, _ := cp.dispatcher.GetCommandType(cmdClassType.(string))

	msg, _ = cp.serializer.Deserialize(message.Value, msg)

	cp.processMessage(msg)

	// commit msg
	cp.reciever.Complete(ctx, message)
}
