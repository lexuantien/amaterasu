package eventsourcing

type IEventSourcedRepository interface {
	Find(string)
	Get(string)
	Save(EventSourced, string)
}
