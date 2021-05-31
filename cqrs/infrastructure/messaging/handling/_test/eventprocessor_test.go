package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/messaging/handling"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/kafkaa"
	"amaterasu/utils"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_publish_event_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"

	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()

	bus := messaging.New_EventBus(sender, serializer)

	err := bus.Publish(context.Background(), messaging.EnvelopeWrap(Event2{
		Id:   utils.NewString(),
		Name: "Tien",
		Div:  12,
	}, messaging.EVENT))

	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func Test_publish_events_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"
	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := messaging.New_Producer(*commandClient)

	serializer := serialization.New_JsonSerializer()
	bus := messaging.New_EventBus(sender, serializer)

	events := []messaging.Envelope{
		messaging.EnvelopeWrap(Event1{
			Name: "Tu",
			Id:   utils.NewString(),
		}, messaging.EVENT),
		messaging.EnvelopeWrap(Event2{
			Id:   utils.NewString(),
			Name: "Ronaldo",
			Div:  1,
		}, messaging.EVENT),
		messaging.EnvelopeWrap(Event3{
			Id:   utils.NewString(),
			Name: "Messi",
			Old:  123,
		}, messaging.EVENT),
	}

	err := bus.Publishes(context.Background(), events...)
	fmt.Println(err)
	assert.Equal(t, nil, err)

}

func Test_subscribe_event_with_1_handler_then_success(t *testing.T) {
	kafkaConfig.Topic = "ni61pj1b-test-event"
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "e107")
	receiver := messaging.New_Consumer(subscriptionClient)
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
	receiver := messaging.New_Consumer(subscriptionClient)
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
