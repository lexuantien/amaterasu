package eventsourcing

type IMementoOriginator interface {
	SaveToMemento() IMemento
}

type IMemento interface {
	Version() int
}
