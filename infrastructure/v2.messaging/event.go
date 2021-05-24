package v2messaging

type IEventHandler interface {
	Handle(interface{}) error
}
