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
		SetTopic(string)
		GetPartition() int32
		SetPartition(int32)
		GetType(string) reflect.Type
		GetVersion() int
		GetId() string
		//
		LoadFromHistory(pastEvents []IVersionedEvent)
		Update(e IVersionedEvent)
		AutoMappingHandles(aggreate interface{}) error
		CreateDefaultValue(string, string, int32, IEventSourced)
	}

	EventSourced struct {
		handlers      map[reflect.Type]func(IVersionedEvent) error // relay event
		pendingEvents []IVersionedEvent                            // pendding event to commit to event_store
		eventTypes    map[string]reflect.Type
		version       int
		id            string
		topic         string
		partition     int32
	}
)

func New_EventSourced() EventSourced {
	return EventSourced{
		handlers:   make(map[reflect.Type]func(IVersionedEvent) error),
		eventTypes: make(map[string]reflect.Type),
		id:         utils.NewUuidString(),
		partition:  -1,
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

func (es *EventSourced) LoadFromHistory(pastEvents []IVersionedEvent) {
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

func (es *EventSourced) GetVersion() int {
	return es.version
}
func (es *EventSourced) GetPartition() int32 {
	return es.partition
}

func (es *EventSourced) SetTopic(topic string) {
	es.topic = topic
}

func (es *EventSourced) SetPartition(partition int32) {
	es.partition = partition
}

func (es *EventSourced) GetType(t string) reflect.Type {
	return es.eventTypes[t]
}

func (es *EventSourced) CreateDefaultValue(id, topic string, partition int32, agg IEventSourced) {
	es.eventTypes = make(map[string]reflect.Type)
	es.handlers = make(map[reflect.Type]func(IVersionedEvent) error)
	es.id = id
	es.topic = topic
	es.partition = partition
	es.AutoMappingHandles(agg)
}

func (es *EventSourced) GetId() string {
	return es.id
}
