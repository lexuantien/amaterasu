package handling

import (
	"errors"
	"leech-service/infrastructure/utils"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"reflect"
)

type CommandDispatcher struct {
	handlers   map[reflect.Type]v2messaging.ICommandHandler
	registries map[string]reflect.Type
}

func New_CommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		handlers:   make(map[reflect.Type]v2messaging.ICommandHandler),
		registries: make(map[string]reflect.Type),
	}
}

// Registers the specified command handler.
func (cd *CommandDispatcher) Register(commandHandler v2messaging.ICommandHandler, commands ...interface{}) error {
	types := make([]_type, len(commands))
	for i, command := range commands {
		t, n := utils.GetTypeName(command)
		if _, ok := cd.handlers[t]; ok {
			return errors.New("the command handled by the received handler already has a registered handler")
		}
		types[i] = _type{t, n}
	}

	// Register this handler for each of he handled types.
	for _, ty := range types {
		cd.handlers[ty.k] = commandHandler
		cd.registries[ty.v] = ty.k
	}

	return nil
}

// Processes the message by calling the registered handler.
func (cd *CommandDispatcher) ProcessMessage(command interface{}) bool {
	t, _ := utils.GetTypeName(command)

	if handler, ok := cd.handlers[t]; ok {
		handler.Handle(command)
		return true
	} else {
		return false
	}
}

func (cd *CommandDispatcher) GetCommandType(name string) (interface{}, error) {
	rawType, ok := cd.registries[name]
	if !ok {
		return nil, errors.New("can't find in registry")
	}
	return reflect.New(rawType).Interface(), nil
}
