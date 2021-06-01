package messaging

type (
	IEventPublisher interface {
		GetEvents() []IEvent
	}

	EventPublisher struct {
		Events []IEvent `gorm:"-"`
	}
)

func (e *EventPublisher) GetEvents() []IEvent {
	return e.Events
}
