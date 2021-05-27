package handling

import (
	kafkaa "amaterasu/cqrs/infrastructure/kafkaa"
	"amaterasu/cqrs/infrastructure/serialization"
	v2messaging "amaterasu/cqrs/infrastructure/v2.messaging"
	"amaterasu/cqrs/infrastructure/v2.messaging/handling"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func Test_send_command_to_kafka_then_success(t *testing.T) {

	commandClient := kafkaa.New_kafkaa(kafkaConfig)
	sender := v2messaging.New_Producer(*commandClient)
	serializer := serialization.New_JsonSerializer()
	bus := v2messaging.New_CommandBus(sender, serializer)

	ok := bus.Send(context.Background(), v2messaging.EnvelopeWrap(Foo1{
		ProductId:   uuid.NewString(),
		Quantity:    uint(10),
		Description: "Bàn phải xuất xứ từ Nhật Bản",
	}, nil, v2messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), v2messaging.EnvelopeWrap(Foo1{
		ProductId:   uuid.NewString(),
		Quantity:    uint(12),
		Description: "Bàn phải xuất xứ từ Lao`",
	}, nil, v2messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), v2messaging.EnvelopeWrap(Foo2{
		OrderId: uuid.NewString(),
	}, nil, v2messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), v2messaging.EnvelopeWrap(Foo3{
		Gender: true,
		FId:    uuid.NewString(),
	}, nil, v2messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), v2messaging.EnvelopeWrap(Foo4{
		T:   "Ronaldo",
		N:   false,
		FId: uuid.NewString(),
	}, nil, v2messaging.COMMAND))
	fmt.Println(ok)

	fmt.Println("done")
}

func Test_process_message_more_handlers_then_success(t *testing.T) {
	config, err := kafkaa.New_KafkaServer(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	receiver := v2messaging.New_Consumer(config)
	serializer := serialization.New_JsonSerializer()

	handler1 := New_CommandHander1()
	handler2 := New_CommandHander2()

	processor := handling.New_CommandProcessor(receiver, serializer)

	processor.Register(handler1)
	processor.Register(handler2)

	wg.Add(1)
	processor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_process_message_with_1_handler_then_success(t *testing.T) {
	subscriptionClient, _ := kafkaa.New_KafkaServer(kafkaConfig, "z107")

	receiver := v2messaging.New_Consumer(subscriptionClient)

	serializer := serialization.New_JsonSerializer()

	orderFooCommandHandler := New_CommandHander1()

	commandProcessor := handling.New_CommandProcessor(receiver, serializer)
	commandProcessor.Register(orderFooCommandHandler)

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_process_message_more_handlers_then_fail(t *testing.T) {
	subscriptionClient, err := kafkaa.New_KafkaServer(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	receiver := v2messaging.New_Consumer(subscriptionClient)

	serializer := serialization.New_JsonSerializer()

	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	commandProcessor := handling.New_CommandProcessor(receiver, serializer)

	commandProcessor.Register(commandHandler1)

	err = commandProcessor.Register(commandHandler2)

	fmt.Println(err)

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}
