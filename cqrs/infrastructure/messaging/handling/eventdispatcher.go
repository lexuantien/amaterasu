package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/utils"
	"errors"
	"reflect"
	"strings"
)

type eventHandlerFunc func(msg messaging.Message) error

//
type EventDispatcher struct {
	eventHandlers     map[reflect.Type]map[messaging.IEventHandler]struct{}
	eventTypes        map[string]reflect.Type
	eventHandlerFuncs map[string][]eventHandlerFunc
}

func New_EventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		eventHandlers:     make(map[reflect.Type]map[messaging.IEventHandler]struct{}),
		eventTypes:        make(map[string]reflect.Type),
		eventHandlerFuncs: make(map[string][]eventHandlerFunc),
	}
}

// Registers the specified command handler.
func (cd *EventDispatcher) Register(eventHandler messaging.ICommandHandler) error {

	handlerType := reflect.TypeOf(eventHandler)
	numHandlerMethods := handlerType.NumMethod() // count all method

	handleFuncCount := 0 // for count handle event func inside eventHandler

	for i := 0; i < numHandlerMethods; i++ {
		// get a method in eventHandler
		handlerMethod := handlerType.Method(i)

		if !strings.HasPrefix(handlerMethod.Name, PREFIX) || // skip method without [PREFIX]
			handlerMethod.Type.NumIn() != NUMIN || // only handle method should have [NUMIN] arguments,
			handlerMethod.Type.NumOut() != NUMOUT { // [NUMOUT] outputs
			continue
		}

		// yeahhh, we have handle event func in command handler, so lucky
		handleFuncCount++

		// ex:
		// func (fh fooEventHandler) handleFoo1(f foo1) error {}
		// func (fh fooEventHandler) handleFoo2(f foo2) error {}
		handlerFunc := func(msg messaging.Message) error {

			response := handlerMethod.Func.Call([]reflect.Value{
				reflect.ValueOf(eventHandler), // fh param
				reflect.ValueOf(msg).Elem(),   // f param
			})

			// can use class if this class is nil
			// ex: func (fh fooCommandHandler) handleFoo1(f foo1) error {}
			// if fh nil, can't call `handleFoo1` func
			if len(response) > 0 && !response[0].IsNil() {
				err := response[0].Interface().(error)
				return err
			}

			return nil
		}

		// 1 because param in golang start at 1, 0 is struct type
		eventType := handlerMethod.Type.In(1) // get event type
		eventTypeName := utils.GetObjType2(eventType)
		if _, ok := cd.eventHandlers[eventType]; !ok {
			cd.eventHandlers[eventType] = make(map[messaging.IEventHandler]struct{})
			cd.eventTypes[eventTypeName] = eventType
		}
		cd.eventHandlers[eventType][eventHandler] = struct{}{}
		cd.eventHandlerFuncs[eventTypeName] = append(cd.eventHandlerFuncs[eventTypeName], handlerFunc)
	}

	if handleFuncCount == 0 {
		return errors.New("must have handleEvent method, please read the document please")
	}

	return nil
}

// Processes the message by calling the registered handler.
func (cd *EventDispatcher) DispatchMessage(msg messaging.Envelope, serializer serialization.ISerializer) error {
	handlerFuncArr, found := cd.eventHandlerFuncs[msg.MsgType]

	if !found {
		return errors.New("not found consumer handler")
	}

	entry, exists := cd.eventTypes[msg.MsgType]

	if !exists {
		return errors.New("not found event type")
	}

	event, err := serializer.Deserialize(msg.Body, entry)
	if err != nil {
		return err
	}

	for _, handlerFunc := range handlerFuncArr {
		errHandlerFunc := handlerFunc(event.(messaging.IEvent))
		if errHandlerFunc != nil {
			return errHandlerFunc
		}
	}

	return nil
}
