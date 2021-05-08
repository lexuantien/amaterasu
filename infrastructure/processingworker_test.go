package infrastructure

import (
	"context"
	"fmt"
	"leech-service/infrastructure/messagebroker"
	"leech-service/infrastructure/messaging"
	"leech-service/infrastructure/messaging/handling"
	"leech-service/infrastructure/serialization"
	"strconv"
	"testing"
)

type OrderCommandHander struct{}

func New_OrderCommandHandler() *OrderCommandHander {
	return &OrderCommandHander{}
}

func (o OrderCommandHander) Handle(command interface{}) error {
	switch command.(type) {

	case *PlaceOrder:
		fmt.Println(1)
	case *ConfirmOrder:
	case *CancelOrder:
	}

	return nil
}

type PlaceOrder struct {
	ProductId   string
	Quantity    uint
	Description string
}

type ConfirmOrder struct {
	OrderId string
}

type CancelOrder struct {
	OrderId string
}

func TestSendCommand(t *testing.T) {

	// connect to kafka server
	client := messagebroker.New_TopicClient(messagebroker.Scram{
		Username: "ni61pj1b",
		Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
		Al256:    true,
	}, []string{
		"glider-01.srvs.cloudkafka.com:9094",
		"glider-02.srvs.cloudkafka.com:9094",
		"glider-03.srvs.cloudkafka.com:9094",
	}, "ni61pj1b--topic-A")

	sender := messaging.New_TopicSender(*client)     // clazz handle send message to kafka
	serializer := serialization.New_JsonSerializer() // serialize data before send

	bus := messaging.New_CommandBus(sender, serializer) // using command to send command

	for i := 0; i < 5; i++ {
		// send to kafka
		ok := bus.Send(context.Background(), messaging.CreateCommand(PlaceOrder{
			ProductId:   "z" + strconv.Itoa(i),
			Quantity:    uint(i),
			Description: "Bàn phải xuất xứ từ Nhật Bản",
		}))

		if ok == nil {
			fmt.Println("Send success")
		} else {
			fmt.Println("Send FAIL")
			panic(ok)
		}
	}
}

func TestProcessCommand(t *testing.T) {

	// connect to kafka server
	client, err := messagebroker.New_SubscriptionClient(messagebroker.Scram{
		Username: "ni61pj1b",
		Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
		Al256:    true,
	}, []string{ // broker server
		"glider-01.srvs.cloudkafka.com:9094",
		"glider-02.srvs.cloudkafka.com:9094",
		"glider-03.srvs.cloudkafka.com:9094",
	},
		"ni61pj1b--topic-A", // coonsume topic
		"z105",              // consumer-group
	)

	if err != nil {
		panic(err)
	}

	// clazz use to fetch message from kafka
	receiver := messaging.New_SubscriptionReciever(client)
	// handle deserialize
	serializer := serialization.New_JsonSerializer()

	// command handler
	orderCommandHandler := New_OrderCommandHandler()

	// command processor
	// 1 command processor handle 1 command handler
	orderCommandProcessor := handling.New_CommandProcessor(receiver, serializer)
	orderCommandProcessor.Register(
		orderCommandHandler, // command handler
		PlaceOrder{},        // command
		ConfirmOrder{},      // command
		CancelOrder{},       // command
	)

	// start fetch data from queue
	orderCommandProcessor.Start()
}
