package v2messaging

import (
	"leech-service/cqrs/infrastructure/uuid"
	"time"
)

type MessageStatus uint

const (
	UPDATE MessageStatus = iota
	INSERT
	DELETE
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
	MsgStatus    MessageStatus
}

func EnvelopeWrap(body interface{}) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         nil,
	}
}

func EnvelopeWrap2(body interface{}, partitionKey int32) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,

		PartitionKey: partitionKey,
		Topic:        nil,
	}
}

func EnvelopeWrap3(body interface{}, topic string) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,
		PartitionKey:  -1,
		Topic:         &topic,
	}
}

func EnvelopeWrap4(body interface{}, partitionKey int32, topic string) Envelope {
	return Envelope{
		Id:            uuid.New(),
		CorrelationId: uuid.New(),
		Body:          body,

		PartitionKey: partitionKey,
		Topic:        &topic,
	}
}
