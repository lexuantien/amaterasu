package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/messaging/handling"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/kafkaa"
	"amaterasu/utils"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_send_command_to_kafka_then_success(t *testing.T) {

	config, _ := kafkaa.New_ProducerConfig(kafkaConfig)

	producer := messaging.New_Producer(*config)
	serializer := serialization.New_JsonSerializer()
	bus := messaging.New_CommandBus(producer, serializer)
	f := Foo1{
		Command:     messaging.Command{Id: "123"},
		ProductId:   utils.NewUuidString(),
		Quantity:    uint(10),
		Description: "Bàn phải xuất xứ từ Nhật Bản",
	}

	ok := bus.Send(context.Background(), messaging.EnvelopeWrap(&f, messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), messaging.EnvelopeWrap(&Foo1{
		ProductId:   utils.NewUuidString(),
		Quantity:    uint(12),
		Description: "Bàn phải xuất xứ từ Lao`",
	}, messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), messaging.EnvelopeWrap(&Foo2{
		OrderId: utils.NewUuidString(),
	}, messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), messaging.EnvelopeWrap(&Foo3{
		Gender: true,
		FId:    utils.NewUuidString(),
	}, messaging.COMMAND))
	fmt.Println(ok)

	ok = bus.Send(context.Background(), messaging.EnvelopeWrap(&Foo4{
		T:   "Ronaldo",
		N:   false,
		FId: utils.NewUuidString(),
	}, messaging.COMMAND))
	fmt.Println(ok)

	fmt.Println("done")
}

func Test_process_message_more_handlers_then_success(t *testing.T) {
	config, err := kafkaa.New_ConsumerConfig(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	consumer := messaging.New_Consumer(config)
	serializer := serialization.New_JsonSerializer()

	handler1 := New_CommandHander1()
	handler2 := New_CommandHander2()

	processor := handling.New_CommandProcessor(consumer, serializer)

	processor.Register(handler1)
	processor.Register(handler2)

	wg.Add(1)
	processor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_process_message_with_1_handler_then_success(t *testing.T) {
	config, _ := kafkaa.New_ConsumerConfig(kafkaConfig, "z107")

	consumer := messaging.New_Consumer(config)

	serializer := serialization.New_JsonSerializer()

	orderFooCommandHandler := New_CommandHander1()

	commandProcessor := handling.New_CommandProcessor(consumer, serializer)
	commandProcessor.Register(orderFooCommandHandler)

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()
	fmt.Println(1)
}

func Test_process_message_more_handlers_then_fail(t *testing.T) {
	config, err := kafkaa.New_ConsumerConfig(kafkaConfig, "z105")

	if err != nil {
		panic(err)
	}

	consumer := messaging.New_Consumer(config)

	serializer := serialization.New_JsonSerializer()

	commandHandler1 := New_CommandHander1()
	commandHandler2 := New_CommandHander2()

	commandProcessor := handling.New_CommandProcessor(consumer, serializer)

	commandProcessor.Register(commandHandler1)

	err = commandProcessor.Register(commandHandler2)

	fmt.Println(err)

	wg.Add(1)
	commandProcessor.Start()
	wg.Wait()

	fmt.Println(1)
}

type (
	// Command struct {
	// 	Id string `json:"Id"`
	// }
	PlaceSeat struct {
		messaging.Command
		SeatId   []string
		Quantity int
	}
)

func Test_1(t *testing.T) {
	cmd := PlaceSeat{
		Command: messaging.Command{
			Id: "uuid",
		},
		SeatId: []string{
			"1", "2",
		},
		Quantity: 10,
	}
	cmdByte := Serialize(cmd)
	/* Send to kafka ...*/

	/* Consume command from kafka ... */
	exportVal := reflect.New(reflect.TypeOf(PlaceSeat{})).Interface()
	fmt.Println(exportVal)
	_ = json.Unmarshal(cmdByte, exportVal)
	fmt.Println(exportVal)
}

func Serialize(payload interface{}) []byte {
	cmdByte, _ := json.Marshal(payload)
	fmt.Println(string(cmdByte))
	return cmdByte
}
