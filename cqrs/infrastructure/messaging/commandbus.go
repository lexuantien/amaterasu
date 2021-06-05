package messaging

import (
	"amaterasu/cqrs/infrastructure/serialization"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ICommandBus interface {
	Send(Envelope) (bool, error)
	Sends(...Envelope) (bool, error)
}

// sends commands to the write model through a command bus
type CommandBus struct {
	// kafka producer
	producer IMessageProducer
	// serializer to serialize data
	serializer serialization.ISerializer
}

// create new command bus
// 	producer kafka producer
// 	ser
// 	*CommandBus
func New_CommandBus(producer IMessageProducer, ser serialization.ISerializer) *CommandBus {
	return &CommandBus{
		producer:   producer,
		serializer: ser,
	}
}

// send message
func (bus CommandBus) Send(ctx context.Context, command Envelope) error {
	message := bus.buildMessage(command)

	return bus.producer.Send(ctx, message) // Send to kafka
}

// send messages
func (bus CommandBus) Sends(ctx context.Context, commands ...Envelope) error {
	for _, command := range commands {
		err := bus.Send(ctx, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bus CommandBus) buildMessage(command Envelope) *kafka.Message {

	message := &kafka.Message{}

	idBytes, _ := bus.serializer.Serialize(command.Id)
	message.Key = idBytes

	cmdBytes, _ := bus.serializer.Serialize(command)
	message.Value = cmdBytes

	message.TopicPartition = kafka.TopicPartition{
		Partition: command.PartitionKey,
		Topic:     command.Topic,
	}

	return message
}
