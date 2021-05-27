package handling

import (
	kafkaa "amaterasu/cqrs/infrastructure/kafkaa"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/cqrs/infrastructure/uuid"
	v2messaging "amaterasu/cqrs/infrastructure/v2.messaging"
	"amaterasu/cqrs/infrastructure/v2.messaging/handling"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_publish_event_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"

	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := v2messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()

	bus := v2messaging.New_EventBus(sender, serializer)

	err := bus.Publish(context.Background(), v2messaging.EnvelopeWrap(Event2{
		Id:   uuid.NewString(),
		Name: "Tien",
		Div:  12,
	}, nil, v2messaging.EVENT))

	assert.Equal(t, nil, err)
}

func Test_publish_events_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"
	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := v2messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()
	bus := v2messaging.New_EventBus(sender, serializer)

	events := []v2messaging.Envelope{
		v2messaging.EnvelopeWrap(Event1{
			Name: "Tu",
			Id:   uuid.NewString(),
		}, nil, v2messaging.EVENT),
		v2messaging.EnvelopeWrap(Event2{
			Id:   uuid.NewString(),
			Name: "Ronaldo",
			Div:  1,
		}, nil, v2messaging.EVENT),
		v2messaging.EnvelopeWrap(Event3{
			Id:   uuid.NewString(),
			Name: "Messi",
			Old:  123,
		}, nil, v2messaging.EVENT),
	}

	err := bus.Publishes(context.Background(), events...)
	fmt.Println(err)
	assert.Equal(t, nil, err)

}

func Test_subscribe_event_with_1_handler_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "e107")
	receiver := v2messaging.New_Consumer(subscriptionClient)
	serializer := serialization.New_JsonSerializer()

	eventHandler := New_EventHandler1()

	eventProcessor := handling.New_EventProcessor(receiver, serializer)
	eventProcessor.Register(eventHandler)

	wg.Add(1)
	eventProcessor.Start()
	wg.Wait()
	fmt.Println("d")
}

func Test_subscribe_events_with_multi_handlers_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "e107")
	receiver := v2messaging.New_Consumer(subscriptionClient)
	serializer := serialization.New_JsonSerializer()

	eventHandler1 := New_EventHandler1()
	eventHandler2 := New_EventHandler2()

	eventProcessor := handling.New_EventProcessor(receiver, serializer)
	eventProcessor.Register(eventHandler1)
	eventProcessor.Register(eventHandler2)

	wg.Add(1)
	eventProcessor.Start()
	wg.Wait()
	fmt.Println("d")
}
