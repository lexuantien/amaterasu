package messaging

type (
	IEventPublisher interface {
		GetEvents() []IEvent
	}

	EventPublisher struct {
		Events []IEvent
	}
)
