package v2messaging

import (
	"amaterasu/cqrs/infrastructure/utils"
	"amaterasu/cqrs/infrastructure/uuid"
	"reflect"
	"time"
)

type MessageAction uint

type Message interface{}

const (
	EVENT MessageAction = iota
	COMMAND
)

// Provides the envelope for an object that will be sent to a bus.
type Envelope struct {

	// Identify
	Id uuid.UUID

	// for tracing and monitoring
	CorrelationId   uuid.UUID
	TraceIdentifier uuid.UUID

	// Wrap the command
	Body interface{}

	// Gets or sets the delay for sending, enqueing or processing the body
	Delay *time.Time
	// Gets or sets the time to live for the message in the queue.
	Time2Live *time.Time

	//
	PartitionKey int32
	Topic        *string
	MsgStatus    *uint

	//

	MsgType   string
	MsgAction MessageAction
}

func EnvelopeWrap(body interface{}, msgStatus *uint, msgAction MessageAction) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         nil,
		MsgStatus:     msgStatus,
		MsgAction:     msgAction,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
	}
}

func EnvelopeWrap2(body interface{}, msgStatus *uint, msgAction MessageAction, partitionKey int32) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  partitionKey,
		Topic:         nil,
		MsgStatus:     msgStatus,
		MsgAction:     msgAction,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
	}
}

func EnvelopeWrap3(body interface{}, msgStatus *uint, msgAction MessageAction, topic string) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         &topic,
		MsgStatus:     msgStatus,
		MsgAction:     msgAction,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
	}
}

func EnvelopeWrap4(body interface{}, msgStatus *uint, msgAction MessageAction, partitionKey int32, topic string) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  partitionKey,
		Topic:         &topic,
		MsgStatus:     msgStatus,
		MsgAction:     msgAction,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
	}
}
