package messaging

import (
	"context"
	"leech-service/infrastructure/serialization"
	"leech-service/infrastructure/uuid"
)

type ICommandBus interface {
	Send(Envelop)
	Sends(...Envelop)
}

type CommandBus struct {
	sender     IMessageSender
	serializer serialization.ISerializer
}

func (bus CommandBus) Send(ctx context.Context, command Envelop) {
	message := bus.buildMessage(command)
	//?? handle command.Delay fail
	bus.sender.Send(ctx, message, command.Time2Live)
}

func (bus CommandBus) Sends(ctx context.Context, commands ...Envelop) {
	for _, command := range commands {
		bus.Send(ctx, command)
	}
}

func (bus CommandBus) buildMessage(command Envelop) []byte {
	message, _ := bus.serializer.Serialize(command.Body)

	if command.Id == uuid.Nil { // create new command id
		command.Id = uuid.New()
	}

	if command.CorrelationId == uuid.Nil {
		command.CorrelationId = uuid.New()
	}

	return message
}
