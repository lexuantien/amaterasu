package v2messaging

type ICommandHandler interface {
	Handle(interface{}) error
}

type ICommandHandlerRegistry interface {
	Register(ICommandHandler, ...interface{}) error
}
