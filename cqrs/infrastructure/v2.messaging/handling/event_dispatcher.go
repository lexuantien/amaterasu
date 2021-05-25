package handling

import (
	"errors"
	"leech-service/cqrs/infrastructure/utils"
	v2messaging "leech-service/cqrs/infrastructure/v2.messaging"
	"reflect"
)

//
type EventDispatcher struct {
	handlers   map[reflect.Type]map[v2messaging.IEventHandler]struct{}
	registries map[string]reflect.Type
}

func New_EventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers:   make(map[reflect.Type]map[v2messaging.IEventHandler]struct{}),
		registries: make(map[string]reflect.Type),
	}
}

// Registers the specified command handler.
func (cd *EventDispatcher) Register(eventHandler v2messaging.ICommandHandler, events ...interface{}) error {
	for _, event := range events {
		eventType, eventTypeName := utils.GetTypeName(event)

		if _, ok := cd.handlers[eventType]; !ok {
			cd.handlers[eventType] = make(map[v2messaging.IEventHandler]struct{})
			cd.registries[eventTypeName] = eventType
		}
		cd.handlers[eventType][eventHandler] = struct{}{}
	}
	return nil
}

// Processes the message by calling the registered handler.
func (cd *EventDispatcher) DispatchMessage(event interface{}) bool {
	eventType, _ := utils.GetTypeName(event)
	if handlers, ok := cd.handlers[eventType]; ok {
		for handler := range handlers {
			handler.Handle(event)
		}
		return true
	}
	return false
}

func (cd *EventDispatcher) GetCommandType(name string) (interface{}, error) {
	rawType, ok := cd.registries[name]
	if !ok {
		return nil, errors.New("can't find in registry")
	}
	return reflect.New(rawType).Interface(), nil
}
