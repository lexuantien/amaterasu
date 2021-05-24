package v2messaging

import (
	"leech-service/infrastructure/uuid"
	"time"
)

type MessageStatus uint

const (
	UPDATE MessageStatus = iota
	INSERT
	DELETE
)

// Command represents an actor intention to alter the state of the system
type Envelop struct {

	// Identify
	Id uuid.UUID

	//
	CorrelationId   uuid.UUID
	TraceIdentifier uuid.UUID

	//
	Body interface{} // Wrap the command

	//
	Delay     *time.Time
	Time2Live *time.Time

	//
	PartitionKey int32
	Topic        *string
	MsgStatus    MessageStatus
}

func EnvelopeWrap(body interface{}) Envelop {
	return Envelop{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         nil,
	}
}

func EnvelopeWrap2(body interface{}, partitionKey int32) Envelop {
	return Envelop{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,

		PartitionKey: partitionKey,
		Topic:        nil,
	}
}

func EnvelopeWrap3(body interface{}, topic string) Envelop {
	return Envelop{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         &topic,
	}
}

func EnvelopeWrap4(body interface{}, partitionKey int32, topic string) Envelop {
	return Envelop{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,

		PartitionKey: partitionKey,
		Topic:        &topic,
	}
}

func (e *Envelop) BodyMsg() interface{} {
	return e.Body
}
