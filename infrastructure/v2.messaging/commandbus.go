package v2messaging

import (
	"context"
	"leech-service/infrastructure/serialization"
	"leech-service/infrastructure/utils"
	"leech-service/infrastructure/uuid"
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
		sen,
		ser,
	}
}

func (bus CommandBus) Send(ctx context.Context, command Envelop) error {
	message, clazzType := bus.buildMessage(command)
	//?? handle command.Delay fail
	/// . . .
	// Send to kafka
	return bus.sender.Send(ctx, message, clazzType)
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

func (bus CommandBus) buildMessage(command Envelop) ([]byte, []byte) {
	message, _ := bus.serializer.Serialize(command.Body)
	_, name := utils.GetTypeName(command.Body)
	clazzType, _ := bus.serializer.Serialize(name)

	if command.Id == uuid.Nil { // create new command id
		command.Id = uuid.New()
	}

	if command.CorrelationId == uuid.Nil {
		command.CorrelationId = uuid.New()
	}

	return message, clazzType
}
