package handling

import (
	"errors"
	"leech-service/infrastructure/utils"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"reflect"
)

//
type EventDispatcher struct {
	handlers   map[reflect.Type][]v2messaging.IEventHandler
	registries map[string]reflect.Type
}

func New_EventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers:   make(map[reflect.Type][]v2messaging.IEventHandler),
		registries: make(map[string]reflect.Type),
	}
}

// Registers the specified command handler.
func (cd *EventDispatcher) Register(eventHandler v2messaging.ICommandHandler, events ...interface{}) error {

	for _, event := range events {
		eventType, eventTypeName := utils.GetTypeName(event)

		if eventHandlers, ok := cd.handlers[eventType]; !ok {
			i, n := 0, len(eventHandlers)
			for ; i < n; i++ {
				if eventHandlers[i] == eventHandler {
					i = n
					break
				}
			}
			if i == n {
				cd.handlers[eventType] = append(cd.handlers[eventType], eventHandler)
			}
		} else {
			cd.handlers[eventType] = append(cd.handlers[eventType], eventHandler)
			cd.registries[eventTypeName] = eventType
		}

	}

	return nil
}

// Processes the message by calling the registered handler.
func (cd *EventDispatcher) DispatchMessage(event interface{}) bool {
	eventType, _ := utils.GetTypeName(event)

	if handlers, ok := cd.handlers[eventType]; ok {
		for i := 0; i < len(handlers); i++ {
			handlers[i].Handle(event)
		}
		return true
	} else {
		return false
	}
}

func (cd *EventDispatcher) GetCommandType(name string) (interface{}, error) {
	rawType, ok := cd.registries[name]
	if !ok {
		return nil, errors.New("can't find in registry")
	}
	return reflect.New(rawType).Interface(), nil
}
