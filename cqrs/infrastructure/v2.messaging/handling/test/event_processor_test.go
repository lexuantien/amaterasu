package handling

import (
	"context"
	"fmt"
	kafkaa "leech-service/cqrs/infrastructure/kafkaa"
	"leech-service/cqrs/infrastructure/serialization"
	v2messaging "leech-service/cqrs/infrastructure/v2.messaging"
	"leech-service/cqrs/infrastructure/v2.messaging/handling"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_publish_event_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test/event"

	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := v2messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()

	bus := v2messaging.New_EventBus(sender, serializer)

	err := bus.Publish(context.Background(), v2messaging.EnvelopeWrap(Event2{"LXT", 11}))

	assert.Equal(t, nil, err)
}

func Test_publish_events_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test/event"
	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := v2messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()
	bus := v2messaging.New_EventBus(sender, serializer)

	events := []v2messaging.Envelope{
		v2messaging.EnvelopeWrap(Event1{"LXT"}),
		v2messaging.EnvelopeWrap(Event2{"LXT", 11}),
		v2messaging.EnvelopeWrap(Event3{"LXT", 15, "x"}),
	}

	err := bus.Publishes(context.Background(), events...)

	assert.Equal(t, nil, err)

}

func Test_subscribe_event_with_1_handler_then_success(t *testing.T) {
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "e107")
	receiver := v2messaging.New_Consumer(subscriptionClient)
	serializer := serialization.New_JsonSerializer()

	eventHandler := New_EventHandler1()

	eventProcessor := handling.New_EventProcessor(receiver, serializer)
	eventProcessor.Register(eventHandler, Event1{}, Event2{})

	wg.Add(1)
	eventProcessor.Start()
	wg.Wait()
	fmt.Println("d")
}

func Test_subscribe_events_with_multi_handlers_then_success(t *testing.T) {
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "e107")
	receiver := v2messaging.New_Consumer(subscriptionClient)
	serializer := serialization.New_JsonSerializer()

	eventHandler1 := New_EventHandler1()
	eventHandler2 := New_EventHandler2()

	eventProcessor := handling.New_EventProcessor(receiver, serializer)
	eventProcessor.Register(eventHandler1, Event1{}, Event2{})
	eventProcessor.Register(eventHandler2, Event3{}, Event2{})

	wg.Add(1)
	eventProcessor.Start()
	wg.Wait()
	fmt.Println("d")
}
