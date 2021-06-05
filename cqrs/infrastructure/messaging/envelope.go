package messaging

import (
	"amaterasu/utils"
	"encoding/json"
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
	Body []byte

	//
	PartitionKey int32
	Topic        *string

	//
	MsgType   string
	MsgAction MessageAction
}

func Wrap(body interface{}, msgAction MessageAction, partitionKey int32, topic *string) Envelope {

	e := Envelope{}
	e.MsgAction = msgAction
	e.PartitionKey = partitionKey
	e.Topic = topic
	e.MsgType = utils.GetObjType2(reflect.TypeOf(body))
	e.Id = utils.NewUuidString()

	switch t := body.(type) {
	case ICommand:
		if t.GetId() == "" {
			t.SetId(e.Id)
		} else {
			e.Id = t.GetId()
		}
	case IEvent:
		if t.GetSourceId() == "" {
			t.SetSourceId(e.Id)
		} else {
			e.Id = t.GetSourceId()
		}
	}

	bodyByte, _ := json.Marshal(body)
	e.Body = bodyByte

	return e
}

func EnvelopeWrap(body interface{}, msgAction MessageAction) Envelope {
	return Wrap(body, msgAction, -1, nil)
}

func EnvelopeWrap2(body interface{}, partitionKey int32, msgAction MessageAction) Envelope {
	return Wrap(body, msgAction, partitionKey, nil)
}

func EnvelopeWrap3(body interface{}, msgAction MessageAction, topic string) Envelope {
	return Wrap(body, msgAction, -1, &topic)
}

func EnvelopeWrap4(body interface{}, msgAction MessageAction, partitionKey int32, topic string) Envelope {
	return Wrap(body, msgAction, partitionKey, &topic)
}
