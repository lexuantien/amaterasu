package eventsourcing

import v2messaging "amaterasu/cqrs/infrastructure/v2.messaging"

type (
	IVersionedEvent interface {
		SetVersion(int)
		GetVersion() int
		v2messaging.IEvent
	}

	VersionedEvent struct {
		version  int
		sourceId string
	}
)

func (e *VersionedEvent) SetVersion(version int) {
	e.version = version
}

func (e *VersionedEvent) GetVersion() int {
	return e.version
}

func (e *VersionedEvent) SetSourceID(id string) {
	e.sourceId = id
}

func (e *VersionedEvent) GetSourceID() string {
	return e.sourceId
}
