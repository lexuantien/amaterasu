package handling

import (
	"errors"
	"leech-service/infrastructure/utils"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"reflect"
)

//
type CommandDispatcher struct {
	handler    v2messaging.ICommandHandler
	registries map[string]reflect.Type
}

func New_CommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		registries: make(map[string]reflect.Type),
	}
}

func (cd *CommandDispatcher) Register(commandHandler v2messaging.ICommandHandler, commands ...interface{}) error {

	cd.handler = commandHandler

	for _, command := range commands {
		rawType, name := utils.GetTypeName(command)
		if _, ok := cd.registries[name]; ok {
			return errors.New("duplicate command")
		}
		cd.registries[name] = rawType
	}

	return nil
}

func (cd *CommandDispatcher) GetCommandType(name string) (interface{}, error) {
	rawType, ok := cd.registries[name]
	if !ok {
		return nil, errors.New("can't find in registry")
	}
	return reflect.New(rawType).Interface(), nil
}

func (cd *CommandDispatcher) Dispatch(command interface{}) {
	cd.handler.Handle(command)
}
