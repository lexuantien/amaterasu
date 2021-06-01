package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"context"
)

/*
In the V2 release, the system used the Windows Azure Service Bus to
deliver all commands to their recipients. This meant that the system
delivered the commands asynchronously. In the V3 release, the MVC
controllers now send their commands synchronously and in-process
in order to improve the response times in the UI by bypassing the
command bus and delivering commands directly to their handlers. In
addition, in the ConferenceProcessor worker role, commands sent to
Order aggregates are sent synchronously in-process using the same
mechanism.
*/
type SyncCommandBus struct {
	bus        messaging.CommandBus
	dispatcher *CommandDispatcher
	serializer serialization.ISerializer
}

func New_SyncCommandBus(bus messaging.CommandBus, serializer serialization.ISerializer) {
	syncCb := &SyncCommandBus{}
	syncCb.bus = bus
	syncCb.serializer = serializer
	syncCb.dispatcher = New_CommandDispatcher()
}

func (scb SyncCommandBus) Register(commandHandler messaging.ICommandHandler) error {
	return scb.dispatcher.Register(commandHandler)
}

func (scb SyncCommandBus) Send(ctx context.Context, command messaging.Envelope) error {
	if !scb.doSend(command) {
		return scb.bus.Send(ctx, command)
	}
	return nil
}

func (scb SyncCommandBus) Sends(ctx context.Context, commands ...messaging.Envelope) error {
	// pending command
	for len(commands) > 0 {
		if scb.doSend(commands[0]) {
			commands = commands[1:]
		} else {
			break
		}
	}

	// Commands were not handled locally. Sending it and all remaining commands through the bus
	if len(commands) > 0 {
		scb.bus.Sends(ctx, commands...)
	}

	return nil
}

func (scb SyncCommandBus) doSend(command messaging.Envelope) bool {
	return scb.dispatcher.ProcessMessage(command, scb.serializer) == nil
}
