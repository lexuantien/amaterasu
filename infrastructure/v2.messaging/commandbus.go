package v2messaging

import (
	"context"
	"leech-service/infrastructure/serialization"
	"leech-service/infrastructure/utils"
	"leech-service/infrastructure/uuid"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ICommandBus interface {
	Send(Envelop) (bool, error)
	Sends(...Envelop) (bool, error)
}

type CommandBus struct {
	sender     IMessageSender
	serializer serialization.ISerializer
}

func New_CommandBus(sen IMessageSender, ser serialization.ISerializer) *CommandBus {
	return &CommandBus{
		sender:     sen,
		serializer: ser,
	}
}

func (bus CommandBus) Send(ctx context.Context, command Envelop) error {
	message := bus.buildMessage(command)

	return bus.sender.Send(ctx, message) // Send to kafka
}

func (bus CommandBus) Sends(ctx context.Context, commands ...Envelop) error {
	for _, command := range commands {
		err := bus.Send(ctx, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bus CommandBus) buildMessage(command Envelop) *kafka.Message {

	message := &kafka.Message{}
	var uid uuid.UUID
	if command.Id == uuid.Nil {
		uid = uuid.New()
	}

	idByte, _ := bus.serializer.Serialize(uid)
	message.Key = idByte

	val, _ := bus.serializer.Serialize(command.Body)
	message.Value = val

	_, name := utils.GetTypeName(command.Body)
	cmdType, _ := bus.serializer.Serialize(name)

	message.Headers = append(message.Headers, kafka.Header{
		Key:   "cmd-type",
		Value: cmdType,
	})

	// TODO handle correlationId
	// TODO handle time2live message
	// TODO handle command delay time

	return message
}
