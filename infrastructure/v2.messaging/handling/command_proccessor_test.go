package handling

import (
	"context"
	"fmt"
	"leech-service/infrastructure/messagebroker"
	"leech-service/infrastructure/serialization"
	v2messaging "leech-service/infrastructure/v2.messaging"
	"strconv"
	"sync"
	"testing"
)

var (
	wg          = sync.WaitGroup{}
	kafkaConfig = messagebroker.KafkaConfig{
		Scr: &messagebroker.Scram{
			Username: "ni61pj1b",
			Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
			Al256:    true,
		},
		Brokers: "glider-01.srvs.cloudkafka.com:9094,glider-02.srvs.cloudkafka.com:9094,glider-03.srvs.cloudkafka.com:9094",
		Topic:   "ni61pj1b--topic-A",
	}
)

func Test_send_command_to_kafka_then_success(t *testing.T) {

	commandClient := messagebroker.New_TopicClient(kafkaConfig)
	sender := v2messaging.New_TopicSender(*commandClient)
	serializer := serialization.New_JsonSerializer()
	bus := v2messaging.New_CommandBus(sender, serializer)

	for i := 0; i < 5; i++ {
		ok := bus.Send(context.Background(), v2messaging.CreateCommand(Foo1{
			ProductId:   "z" + strconv.Itoa(i),
			Quantity:    uint(i),
			Description: "Bàn phải xuất xứ từ Nhật FOO Bản",
		}))

		if ok == nil {
			fmt.Println("Send success")
		} else {
			fmt.Println("Send FAIL")
			panic(ok)
		}
	}
	fmt.Println("done")
}

func Test_process_message_with_1_handler_then_success(t *testing.T) {
	subscriptionClient, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z107")

	if err != nil {
		panic(err)
	}

	receiver := v2messaging.New_SubscriptionReciever(subscriptionClient)

	serializer := serialization.New_JsonSerializer()

	orderFooCommandHandler := New_CommandHander1()

	commandProcessor := New_CommandProcessor(receiver, serializer)
	commandProcessor.Register(orderFooCommandHandler, Foo1{}, Foo2{}, Foo3{})

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_process_message_more_handlers_then_fail(t *testing.T) {
	subscriptionClient, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	receiver := v2messaging.New_SubscriptionReciever(subscriptionClient)

	serializer := serialization.New_JsonSerializer()

	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	commandProcessor := New_CommandProcessor(receiver, serializer)

	commandProcessor.Register(
		commandHandler1,
		Foo1{},
		Foo2{},
		Foo3{},
	)

	err = commandProcessor.Register(
		commandHandler2,
		Foo1{}, // this command will error because duplicate
	)

	fmt.Println(err)

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}

func Test_process_message_more_handlers_then_success(t *testing.T) {
	subscriptionClient, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	receiver := v2messaging.New_SubscriptionReciever(subscriptionClient)

	serializer := serialization.New_JsonSerializer()

	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	commandProcessor := New_CommandProcessor(receiver, serializer)

	commandProcessor.Register(
		commandHandler1,
		Foo1{},
		Foo2{},
		Foo3{},
	)

	err = commandProcessor.Register(
		commandHandler2,
		Foo4{}, // this command will success
	)

	fmt.Println(err) // nil

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}
