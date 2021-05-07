package infrastructure

import (
	"leech-service/infrastructure/messagebroker"
	"leech-service/infrastructure/messaging"
	"leech-service/infrastructure/messaging/handling"
	"leech-service/infrastructure/serialization"
	"testing"
)

type OrderCommandHander struct{}

func New_OrderCommandHandler() *OrderCommandHander {
	return &OrderCommandHander{}
}

func (o OrderCommandHander) Handle(command messaging.ICommand) error {
	switch command.Command().(type) {

	case *PlaceOrder:
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

func TestProcessCommand(t *testing.T) {

	client, err := messagebroker.New_SubscriptionClient(messagebroker.Scram{
		Username: "ni61pj1b",
		Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
		Al256:    true,
	}, []string{
		"glider-01.srvs.cloudkafka.com:9094",
		"glider-02.srvs.cloudkafka.com:9094",
		"glider-03.srvs.cloudkafka.com:9094",
	}, "ni61pj1b-topic-A", "z105")

	if err != nil {
		panic(err)
	}

	receiver := messaging.New_SubscriptionReciever(client)
	serializer := serialization.New_JsonSerializer()

	// command processor
	orderCommandProcessor := handling.New_CommandProcessor(receiver, serializer)

	// command handler
	orderCommandHandler := New_OrderCommandHandler()
	orderCommandProcessor.Register(
		orderCommandHandler,
		&PlaceOrder{},
		&ConfirmOrder{},
		&CancelOrder{},
	)

	// start fetch data from queue
	orderCommandProcessor.Start()
}
