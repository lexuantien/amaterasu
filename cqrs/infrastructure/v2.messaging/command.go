package v2messaging

type ICommandHandler interface {
}

type ICommandHandlerRegistry interface {
	Register(ICommandHandler) error
}

type ICommand interface {
	Id() string
}
