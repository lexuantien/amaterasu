package v2messaging

type ICommandHandler interface {
	Handle(interface{}) error
}
