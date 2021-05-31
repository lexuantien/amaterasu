package messaging

type (
	IEventHandler interface{}

	IEventHandlerRegistry interface {
		Register(IEventHandler) error
	}

	IEvent interface {
		SetSourceID(string)
		GetSourceID() string
	}

	Event struct {
		SourceID string
	}
)

func (e *Event) GetSourceID() string {
	return e.SourceID
}

func (e *Event) SetSourceID(sourceID string) {
	e.SourceID = sourceID
}
