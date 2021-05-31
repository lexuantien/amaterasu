package messaging

import (
	"amaterasu/utils"
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
	Id string

	// for tracing and monitoring
	CorrelationId   string
	TraceIdentifier string

	// Wrap the command
	Body interface{}

	// Gets or sets the delay for sending, enqueing or processing the body
	Delay *time.Time
	// Gets or sets the time to live for the message in the queue.
	Time2Live *time.Time

	//
	PartitionKey int32
	Topic        *string

	//

	MsgType   string
	MsgAction MessageAction
}

func EnvelopeWrap(body interface{}, msgAction MessageAction) Envelope {
	return Envelope{
		Id:            utils.NewString(),
		CorrelationId: utils.NewString(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         nil,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
		MsgAction:     msgAction,
	}
}

func EnvelopeWrap2(body interface{}, partitionKey int32, msgAction MessageAction) Envelope {
	return Envelope{
		Id:            utils.NewString(),
		CorrelationId: utils.NewString(),
		Body:          body,
		PartitionKey:  partitionKey,
		Topic:         nil,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
		MsgAction:     msgAction,
	}
}

func EnvelopeWrap3(body interface{}, msgAction MessageAction, topic string) Envelope {
	return Envelope{
		Id:            utils.NewString(),
		CorrelationId: utils.NewString(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         &topic,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
		MsgAction:     msgAction,
	}
}

func EnvelopeWrap4(body interface{}, msgAction MessageAction, partitionKey int32, topic string) Envelope {
	return Envelope{
		Id:            utils.NewString(),
		CorrelationId: utils.NewString(),
		Body:          body,
		PartitionKey:  partitionKey,
		Topic:         &topic,
		MsgType:       utils.GetTypeName2(reflect.TypeOf(body)),
		MsgAction:     msgAction,
	}
}
