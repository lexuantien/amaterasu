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

func (cp *CommandProcessor) processMessage(command messaging.ICommand) error {
	return cp.dispatcher.Dispatch(command)
}

func (cp *CommandProcessor) Start() {
	// lock
	cp.lock.Lock()
	cp.reciever.Start(cp.onMessageReceived)
	cp.started = true
	cp.lock.Unlock()
}

func (mp *CommandProcessor) Stop() {

}

// recieve message from message queue
func (cp *CommandProcessor) onMessageReceived(message kafka.Message) {

	payload, err := cp.serializer.Deserialize(message.Value)
	if err != nil {
		// add to dead letter
		return
	}

	// handle payload here
	cp.processMessage(payload.(messaging.ICommand))

	// commit msg
	err = cp.reciever.CommitMessage(context.Background(), message)
	if err != nil {
		// add to dead letter
		return
	}

}
