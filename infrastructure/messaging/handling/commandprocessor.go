package handling

import (
	"context"
	"leech-service/infrastructure/messaging"
	"leech-service/infrastructure/serialization"
	"sync"

	"github.com/segmentio/kafka-go"
)

type IMessageProcesser interface {
	Start()
	Stop()
}

type CommandProcessor struct {
	lock       sync.Mutex
	reciever   messaging.IMessageReceiver
	serializer serialization.ISerializer
	started    bool
	cancelFunc context.CancelFunc
	dispatcher *CommandDispatcher
}

// :Create
func New_CommandProcessor(r messaging.IMessageReceiver, s serialization.ISerializer) *CommandProcessor {
	return &CommandProcessor{
		reciever:   r,
		serializer: s,
		dispatcher: New_CommandDispatcher(),
	}
}

func (cp *CommandProcessor) Register(commandHandler messaging.ICommandHandler, commands ...interface{}) error {
	return cp.dispatcher.Register(commandHandler, commands...)
}

func (cp *CommandProcessor) processMessage(command interface{}) {
	cp.dispatcher.Dispatch(command)
}

func (cp *CommandProcessor) Start() {
	// lock
	cp.lock.Lock()
	ctx, cancelFunc := context.WithCancel(context.Background())
	cp.cancelFunc = cancelFunc
	cp.reciever.Start(ctx, cp.onMessageReceived)
	cp.started = true
	cp.lock.Unlock()
}

func (cp *CommandProcessor) Stop() {
	cp.lock.Lock()
	cp.reciever.Stop(cp.cancelFunc)
	cp.started = false
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
