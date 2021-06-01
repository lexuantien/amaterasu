package messaging

type (
	IEventHandler interface{}

	IEventHandlerRegistry interface {
		Register(IEventHandler) error
	}

	IEvent interface {
		SetSourceId(string)
		GetSourceId() string
	}

	Event struct {
		SourceId string
	}
)

func (e *Event) GetSourceId() string {
	return e.SourceId
}

func (e *Event) SetSourceId(SourceId string) {
	e.SourceId = SourceId
}
