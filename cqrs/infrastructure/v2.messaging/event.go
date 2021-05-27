package v2messaging

type IEventHandler interface{}

type IEventHandlerRegistry interface {
	Register(IEventHandler) error
}

type IEvent interface {
	SourceId() string
}
