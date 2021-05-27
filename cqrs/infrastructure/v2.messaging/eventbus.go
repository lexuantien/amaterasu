package v2messaging

import (
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/cqrs/infrastructure/uuid"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IEventBus interface {
	Publish(event Envelope) (bool, error)
	Publishes(events ...Envelope) (bool, error)
}

// An event bus that sends serialized object payloads through IMessageProducer
type EventBus struct {
	producer   IMessageProducer
	serializer serialization.ISerializer
}

func New_EventBus(sen IMessageProducer, ser serialization.ISerializer) *EventBus {
	return &EventBus{
		producer:   sen,
		serializer: ser,
	}
}

func (bus EventBus) Publish(ctx context.Context, event Envelope) error {
	message := bus.buildMessage(event)

	return bus.producer.Send(ctx, message) // Send to kafka
}

func (bus EventBus) Publishes(ctx context.Context, events ...Envelope) error {
	for _, event := range events {
		err := bus.Publish(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bus EventBus) buildMessage(event Envelope) *kafka.Message {
	message := &kafka.Message{}

	var uid uuid.UUID = event.Id
	if event.Id == uuid.Nil {
		uid = uuid.New()
	}

	idBytes, _ := bus.serializer.Serialize(uid)
	message.Key = idBytes

	cmdBytes, _ := bus.serializer.Serialize(event)
	message.Value = cmdBytes

	// TODO handle correlationId
	// TODO handle time2live message
	// TODO handle command delay time

	return message
}
