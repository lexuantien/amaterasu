package eventsourcing

import (
	"amaterasu/utils"
	"errors"
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
		//
		LoadFrom(pastEvents []IVersionedEvent)
		Update(e IVersionedEvent)
		AutoMappingHandles(aggreate interface{}) error
		CreateDefaultValue(string, string, int, IEventSourced)
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

func New_EventSourced(partition int, topic string) EventSourced {
	return EventSourced{
		handlers:   make(map[reflect.Type]func(IVersionedEvent) error),
		eventTypes: make(map[string]reflect.Type),
		topic:      topic,
		partition:  partition,
		id:         utils.NewUuidString(),
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
		eventTypeName := utils.GetObjType2(eventType)
		es.eventTypes[eventTypeName] = eventType
	}

	return nil
}

func (es *EventSourced) Update(e IVersionedEvent) {
	e.SetSourceId(es.id)
	e.SetVersion(es.version + 1)
	if onEventFunc, found := es.handlers[reflect.TypeOf(e).Elem()]; found {
		onEventFunc(e)
	}
	es.version = e.GetVersion()
	es.pendingEvents = append(es.pendingEvents, e)
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

func (es *EventSourced) CreateDefaultValue(id, topic string, partition int, agg IEventSourced) {
	es.eventTypes = make(map[string]reflect.Type)
	es.handlers = make(map[reflect.Type]func(IVersionedEvent) error)
	es.topic = topic
	es.partition = partition
	es.id = id
	es.AutoMappingHandles(agg)
}
