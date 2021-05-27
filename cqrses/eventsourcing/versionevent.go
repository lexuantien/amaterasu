package eventsourcing

type IVersionedEvent interface {
	Version() int
}

type VersionedEvent struct {
	version  int
	sourceId string
}

func (ve VersionedEvent) SourceId() string { // IEvent interface
	return ve.sourceId
}
