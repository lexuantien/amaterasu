package messaging

type (
	IEventPublisher interface {
		GetEvents() []IEvent
		// Convert2EventByte()
	}

	EventPublisher struct {
		Events []IEvent `gorm:"-"`
		// EventByte []byte   `gorm:"column:events;type:JSON"`
	}
)

func (e *EventPublisher) GetEvents() []IEvent {
	return e.Events
}

// func (e *EventPublisher) Convert2EventByte() {
// 	eventByte, _ := json.Marshal(e.Events)
// 	e.EventByte = eventByte
// }
