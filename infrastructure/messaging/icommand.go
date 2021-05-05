package messaging

import (
	"reflect"
)

// Command represents an actor intention to alter the state of the system
type Command struct {
	Id   string // Identify
	Type string
	Body interface{} // Wrap the command
}

func CreateCommand(body interface{}) Command {
	t := reflect.TypeOf(body)
	return Command{
		Id:   "generate later",
		Type: t.String(),
		Body: body,
	}
}

//
type icommandhandler interface {
	Add(commands ...Command) error
	Handle(command Command) error
}
