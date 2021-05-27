package eventsourcing

import "reflect"

type VersionedEvent struct{}

// https://github.com/microsoftarchive/cqrs-journey/blob/master/docs/Reference_04_DeepDive.markdown
// a set of classes that define an Order aggregate
type EventSourced struct {
	//
	handlers      map[reflect.Type]func(VersionedEvent)
	pendingEvents []VersionedEvent

	//
	id      string
	version int
}

func New_EventSourced() *EventSourced {
	return &EventSourced{
		version:  -1,
		handlers: make(map[reflect.Type]func(VersionedEvent)),
	}
}
