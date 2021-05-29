package eventsourcing

import (
	"amaterasu/cqrs/infrastructure/utils"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	HANDLE_EVENTSOURCED_PREFIX = "On"
	NUMIN                      = 2
	NUMOUT                     = 1
)

type (
	IEventSourced interface {
		Events() []IVersionedEvent
		GetTopic() string
		GetPartition() int
		GetType(string) reflect.Type
	}

	EventSourced struct {
		handlers      map[reflect.Type]func(IVersionedEvent) error // relay event
		pendingEvents []IVersionedEvent                            // pendding event to commit to event_store
		eventTypes    map[string]reflect.Type
		version       int
		id            string
		topic         string
		partition     int
	}
)

func New_EventSourced() EventSourced {
	return EventSourced{
		handlers:   make(map[reflect.Type]func(IVersionedEvent) error),
		eventTypes: make(map[string]reflect.Type),
	}
}

func (es *EventSourced) AutoMappingHandles(aggreate interface{}) error {
	aggType := reflect.TypeOf(aggreate)
	numAggMethod := aggType.NumMethod()

	if numAggMethod == 0 {
		return errors.New("must have onEvent method, please read the document please")
	}

	eventFuncCount := 0 // for count handle event func

	for i := 0; i < numAggMethod; i++ {
		// get a method in
		aggMethod := aggType.Method(i)
		if !strings.HasPrefix(aggMethod.Name, HANDLE_EVENTSOURCED_PREFIX) || // skip method without [PREFIX]
			aggMethod.Type.NumIn() != NUMIN || // only handle method should have [NUMIN] arguments,
			aggMethod.Type.NumOut() != NUMOUT { // [NUMOUT] outputs
			continue
		}

		eventFuncCount++

		onEventFunc := func(e IVersionedEvent) error {

			response := aggMethod.Func.Call([]reflect.Value{
				reflect.ValueOf(aggreate), // fh param
				reflect.ValueOf(e).Elem(), // f param
			})

			// can use class if this class is nil
			if len(response) > 0 && !response[0].IsNil() {
				err := response[0].Interface().(error)
				return err
			}

			return nil
		}

		eventType := aggMethod.Type.In(1)

		if _, found := es.handlers[eventType]; found {
			return errors.New("can not have 2 method handling 1 event")
		}

		es.handlers[eventType] = onEventFunc
		eventTypeName := utils.GetTypeName2(eventType)
		es.eventTypes[eventTypeName] = eventType
	}

	return nil
}

func (es *EventSourced) Update(e IVersionedEvent) {
	// e.SourceId = this.Id;
	// e.Version = this.version + 1;
	// this.handlers[e.GetType()].Invoke(e);
	// this.version = e.Version;
	// this.pendingEvents.Add(e);
	e.SetSourceID(es.id)
	e.SetVersion(es.version + 1)
	if onEventFunc, found := es.handlers[reflect.TypeOf(e).Elem()]; found {
		onEventFunc(e)
	}
	es.version = e.GetVersion()
	es.pendingEvents = append(es.pendingEvents, e)
	fmt.Println(es.pendingEvents[0])
}

func (es *EventSourced) LoadFrom(pastEvents []IVersionedEvent) {
	// foreach (var e in pastEvents) {
	// 	this.handlers[e.GetType()].Invoke(e);
	// 	this.version = e.Version;
	// }

	for _, e := range pastEvents {
		if onEventFunc, found := es.handlers[reflect.TypeOf(e).Elem()]; found {
			onEventFunc(e)
		}
		es.version = e.GetVersion()
	}
}

func (es *EventSourced) Events() []IVersionedEvent {
	return es.pendingEvents
}

func (es *EventSourced) GetTopic() string {
	return es.topic
}

func (es *EventSourced) GetPartition() int {
	return es.partition
}

func (es *EventSourced) GetType(t string) reflect.Type {
	return es.eventTypes[t]
}
