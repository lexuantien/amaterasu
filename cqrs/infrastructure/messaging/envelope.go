package messaging

import (
	"amaterasu/utils"
	"reflect"
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

	// Wrap the command/event
	Body interface{}

	//
	PartitionKey int32
	Topic        *string

	//
	MsgType   string
	MsgAction MessageAction
}

func EnvelopeWrap(body interface{}, msgAction MessageAction) Envelope {
	return Envelope{
		Id:           utils.NewUuidString(),
		Body:         body,
		PartitionKey: -1,
		Topic:        nil,
		MsgType:      utils.GetObjType2(reflect.TypeOf(body)),
		MsgAction:    msgAction,
	}
}

func EnvelopeWrap2(body interface{}, partitionKey int32, msgAction MessageAction) Envelope {
	return Envelope{
		Id:           utils.NewUuidString(),
		Body:         body,
		PartitionKey: partitionKey,
		Topic:        nil,
		MsgType:      utils.GetObjType2(reflect.TypeOf(body)),
		MsgAction:    msgAction,
	}
}

func EnvelopeWrap3(body interface{}, msgAction MessageAction, topic string) Envelope {
	return Envelope{
		Id:           utils.NewUuidString(),
		Body:         body,
		PartitionKey: -1,
		Topic:        &topic,
		MsgType:      utils.GetObjType2(reflect.TypeOf(body)),
		MsgAction:    msgAction,
	}
}

func EnvelopeWrap4(body interface{}, msgAction MessageAction, partitionKey int32, topic string) Envelope {
	return Envelope{
		Id:           utils.NewUuidString(),
		Body:         body,
		PartitionKey: partitionKey,
		Topic:        &topic,
		MsgType:      utils.GetObjType2(reflect.TypeOf(body)),
		MsgAction:    msgAction,
	}
}
