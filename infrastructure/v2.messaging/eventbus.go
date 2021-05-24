package v2messaging

import (
	"context"
	"leech-service/infrastructure/serialization"
	"leech-service/infrastructure/utils"
	"leech-service/infrastructure/uuid"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IEventBus interface {
	Publish(event Envelop) (bool, error)
	Publishes(events ...Envelop) (bool, error)
}

// An event bus that sends serialized object payloads through IMessageSender
type EventBus struct {
	sender     IMessageSender
	serializer serialization.ISerializer
}

func New_EventBus(sen IMessageSender, ser serialization.ISerializer) *EventBus {
	return &EventBus{
		sender:     sen,
		serializer: ser,
	}
}

func (bus EventBus) Publish(ctx context.Context, event Envelop) error {
	message := bus.buildMessage(event)

	return bus.sender.Send(ctx, message) // Send to kafka
}

func (bus EventBus) Publishes(ctx context.Context, events ...Envelop) error {
	for _, event := range events {
		err := bus.Publish(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bus EventBus) buildMessage(event Envelop) *kafka.Message {
	message := &kafka.Message{}

	var uid uuid.UUID
	if event.Id == uuid.Nil {
		uid = uuid.New()
	}

	idByte, _ := bus.serializer.Serialize(uid)
	message.Key = idByte

	val, _ := bus.serializer.Serialize(event.Body)
	message.Value = val

	_, name := utils.GetTypeName(event.Body)
	evtType, _ := bus.serializer.Serialize(name)

	message.Headers = append(message.Headers, kafka.Header{
		Key:   "evt-type",
		Value: evtType,
	})

	message.TopicPartition = kafka.TopicPartition{
		Topic:     event.Topic,
		Partition: event.PartitionKey,
	}

	// TODO handle correlationId
	// TODO handle time2live message
	// TODO handle command delay time

	return message
}
