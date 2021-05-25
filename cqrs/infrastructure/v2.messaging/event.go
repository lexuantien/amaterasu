package v2messaging

type IEventHandler interface {
	Handle(interface{}) error
}

type IEventHandlerRegistry interface {
	Register(IEventHandler, ...interface{}) error
}
