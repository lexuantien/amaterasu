package handling

import (
	"errors"
	"leech-service/infrastructure/messaging"
	"reflect"
)

type CommandDispatcher struct {
	handlers map[string]messaging.ICommandHandler
}

func New_CommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		handlers: make(map[string]messaging.ICommandHandler),
	}
}

func (cd *CommandDispatcher) Register(commandHandler messaging.ICommandHandler, commands ...interface{}) error {
	for _, command := range commands {
		typeName := reflect.TypeOf(command).Elem().Name()
		if _, ok := cd.handlers[typeName]; ok {
			return errors.New("duplicate command")
		}
		cd.handlers[typeName] = commandHandler
	}

	return nil
}

func (cd *CommandDispatcher) Dispatch(command messaging.ICommand) error {
	if commandHandler, ok := cd.handlers[command.CommandType()]; ok {
		commandHandler.Handle(command)
	}

	return errors.New("command handler not found")
}
