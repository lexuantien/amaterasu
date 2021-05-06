package messaging

import (
	"leech-service/infrastructure/uuid"
	"reflect"
	"time"
)

// Command represents an actor intention to alter the state of the system
type Command struct {
	Id            uuid.UUID // Identify
	CorrelationId uuid.UUID
	Type          string
	Body          interface{} // Wrap the command
	Delay         *time.Time
	Time2Live     *time.Time
}

func CreateCommand(body interface{}) Command {
	t := reflect.TypeOf(body)
	return Command{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Type:          t.String(),
		Body:          body,
	}
}
