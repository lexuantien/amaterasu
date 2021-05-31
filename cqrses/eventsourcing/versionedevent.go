package eventsourcing

import "amaterasu/cqrs/infrastructure/messaging"

type (
	IVersionedEvent interface {
		messaging.IEvent
		SetVersion(int)
		GetVersion() int
	}

	VersionedEvent struct {
		messaging.Event
		Version int
	}
)

func (e *VersionedEvent) SetVersion(version int) {
	e.Version = version
}

func (e *VersionedEvent) GetVersion() int {
	return e.Version
}
