package eventsourcing

import (
	"reflect"
)

type IEventSourced interface {
	Id() string

	Version() int

	Events() []VersionedEvent
}

// https://github.com/microsoftarchive/cqrs-journey/blob/master/docs/Reference_04_DeepDive.markdown
// a set of classes that define an Order aggregate
type EventSourced struct {
	//
	handlers      map[reflect.Type]OnEventExcuted
	pendingEvents []VersionedEvent

	//
	id      string
	version int
}

func New_EventSourced() *EventSourced {
	return &EventSourced{
		version:  -1,
		handlers: make(map[reflect.Type]OnEventExcuted),
	}
}

type OnEventExcuted func(VersionedEvent)

func (es *EventSourced) Handles(event EventSourced, onEventExcuted OnEventExcuted) {
	es.handlers[reflect.TypeOf(event)] = onEventExcuted
}

func (es *EventSourced) LoadFrom(pastEvents []VersionedEvent) {
	for _, e := range pastEvents {
		es.handlers[reflect.TypeOf(e)](e) // invoke function
		es.version = e.version
	}
}

func (es *EventSourced) Update(e VersionedEvent) {
	e.sourceId = es.id
	e.version = es.version + 1
	es.handlers[reflect.TypeOf(e)](e)
	es.version = e.version
	es.pendingEvents = append(es.pendingEvents, e)
}

// IEventSourced interface
func (es *EventSourced) Events() []VersionedEvent {
	return es.pendingEvents
}

func (es *EventSourced) Id() string {
	return es.id
}

func (es *EventSourced) Version() int {
	return es.version
}
