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
		Brokers: []string{ // broker server
			"glider-01.srvs.cloudkafka.com:9094",
			"glider-02.srvs.cloudkafka.com:9094",
			"glider-03.srvs.cloudkafka.com:9094",
		},
		Topic: "ni61pj1b--topic-A",
	}
)

func Test_send_command_to_kafka_then_success(t *testing.T) {

	client := messagebroker.New_TopicClient(kafkaConfig)
	sender := v2messaging.New_TopicSender(*client)   // clazz handle send message to kafka
	serializer := serialization.New_JsonSerializer() // serialize data before send

	bus := v2messaging.New_CommandBus(sender, serializer) // using command to send command

	for i := 0; i < 5; i++ {
		// send to kafka
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

func Test_register_1_command_handler_then_success(t *testing.T) {
	client, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	// class use to fetch message from kafka
	receiver := v2messaging.New_SubscriptionReciever(client)

	// handle deserialize
	serializer := serialization.New_JsonSerializer()

	// command handler
	orderFooCommandHandler := New_CommandHander1()

	// create command processor
	commandProcessor := New_CommandProcessor(receiver, serializer)

	// register handler to processor
	commandProcessor.Register(
		orderFooCommandHandler, // command handler
		Foo1{},                 // command
		Foo2{},                 // command
		Foo3{},                 // command
	)

	wg.Add(1)
	// start fetch data from queue
	commandProcessor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_register_more_handlers_then_fail(t *testing.T) {
	client, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	// class use to fetch message from kafka
	receiver := v2messaging.New_SubscriptionReciever(client)

	// handle deserialize
	serializer := serialization.New_JsonSerializer()

	// command handlers
	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	// create command processor
	commandProcessor := New_CommandProcessor(receiver, serializer)

	// register handler to processor
	commandProcessor.Register(
		commandHandler1, // command handler
		Foo1{},          // command
		Foo2{},          // command
		Foo3{},          // command
	)

	err = commandProcessor.Register(
		commandHandler2, // command handler
		Foo1{},          // this command will error because duplicate
	)

	fmt.Println(err)

	wg.Add(1)
	// start fetch data from queue
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}

func Test_register_more_handlers_then_success(t *testing.T) {
	client, err := messagebroker.New_SubscriptionClient(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	// class use to fetch message from kafka
	receiver := v2messaging.New_SubscriptionReciever(client)

	// handle deserialize
	serializer := serialization.New_JsonSerializer()

	// command handlers
	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	// create command processor
	commandProcessor := New_CommandProcessor(receiver, serializer)

	// register handler to processor
	commandProcessor.Register(
		commandHandler1, // command handler
		Foo1{},          // command
		Foo2{},          // command
		Foo3{},          // command
	)

	err = commandProcessor.Register(
		commandHandler2, // command handler
		Foo4{},          // this command will success
	)

	fmt.Println(err) // nil

	wg.Add(1)
	// start fetch data from queue
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}
