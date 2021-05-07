package messaging

import (
	"leech-service/infrastructure/uuid"
	"reflect"
	"time"
)

type ICommand interface {
	Command() interface{}
	CommandType() string
}

// Command represents an actor intention to alter the state of the system
type Envelop struct {
	Id            uuid.UUID // Identify
	CorrelationId uuid.UUID
	Body          interface{} // Wrap the command
	Delay         *time.Time
	Time2Live     *time.Time
}

func CreateCommand(command interface{}) *Envelop {
	return &Envelop{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          command,
	}
}

func (e *Envelop) Command() interface{} {
	return e.Body
}

func (e *Envelop) CommandType() string {
	return reflect.TypeOf(e.Body).Elem().Name()
}

//

type ICommandHandler interface {
	Handle(ICommand) error
}
