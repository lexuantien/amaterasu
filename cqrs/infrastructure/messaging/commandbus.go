package messaging

import (
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/utils"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ICommandBus interface {
	Send(Envelope) (bool, error)
	Sends(...Envelope) (bool, error)
}

type CommandBus struct {
	producer   IMessageProducer
	serializer serialization.ISerializer
}

func New_CommandBus(producer IMessageProducer, ser serialization.ISerializer) *CommandBus {
	return &CommandBus{
		producer:   producer,
		serializer: ser,
	}
}

func (bus CommandBus) Send(ctx context.Context, command Envelope) error {
	message := bus.buildMessage(command)

	return bus.producer.Send(ctx, message) // Send to kafka
}

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

	var uid string = command.Id
	if command.Id == "" {
		uid = utils.NewString()
	}

	idBytes, _ := bus.serializer.Serialize(uid)
	message.Key = idBytes

	cmdBytes, _ := bus.serializer.Serialize(command)
	message.Value = cmdBytes

	// TODO handle correlationId
	// TODO handle time2live message
	// TODO handle command delay time

	return message
}
