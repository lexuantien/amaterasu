package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"context"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Processes incoming commands from the bus and routes them to the appropriate
type CommandProcessor struct {

	// kafka consumer
	consumer messaging.IMessageConsumer
	// serialize kafka message to json - grpc - html ...
	serializer serialization.ISerializer

	// ignore the second request
	started bool

	// uses a peek/lock technique to retrieve a message from a consumer
	lockObject sync.Mutex
	cancelFunc context.CancelFunc

	// dispatch command to specific command handler
	dispatcher *CommandDispatcher
}

// Create new command processor
// Initializes a new instance of the CommandProcessor
// @param r Kafka receive maeesage
// @param s The serializer to use for the message body.
func New_CommandProcessor(r messaging.IMessageConsumer, s serialization.ISerializer) *CommandProcessor {
	return &CommandProcessor{
		consumer:   r,
		serializer: s,
		dispatcher: New_CommandDispatcher(),
	}
}

// Registers the specified command handler.
func (cp *CommandProcessor) Register(commandHandler messaging.ICommandHandler) error {
	return cp.dispatcher.Register(commandHandler)
}

// Processes the message by calling the registered handler.
func (cp *CommandProcessor) processMessage(msg messaging.Envelope) bool {
	return cp.dispatcher.ProcessMessage(msg, cp.serializer) == nil
}

// start consume message in message queue
func (cp *CommandProcessor) Start() {
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

// stop consume message in message queue
func (cp *CommandProcessor) Stop() {
	cp.lockObject.Lock()
	if cp.started {
		cp.consumer.Stop(cp.cancelFunc)
		cp.started = false
	}
	cp.lockObject.Unlock()
}

// recieve message from message queue
func (cp *CommandProcessor) onMessageReceived(ctx context.Context, message *kafka.Message) error {
	return OnMessageReceivedHandler(ctx, message, cp.processMessage, cp.consumer.Complete)
}
