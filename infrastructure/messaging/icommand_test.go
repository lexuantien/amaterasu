package messaging

import (
	"testing"
)

type OrderCommandHander struct{}

//
type PlaceOrder struct {
	ProductId   string
	Quantity    uint
	Description string
}

//
type ConfirmOrder struct {
	OrderId string
}

//
type CancelOrder struct {
	OrderId string
}

//
func TestCreateCommandHandler(t *testing.T) {
	// create command
	placeOrder := CreateCommand(PlaceOrder{ProductId: "z109", Quantity: 10, Description: "this is a book"})
	confirmOrder := CreateCommand(ConfirmOrder{OrderId: "z109"})
	cancelOrder := CreateCommand(CancelOrder{OrderId: "z169"})
	// create command handler
	orderCommandHander := OrderCommandHander{}
	orderCommandHander.Add(placeOrder, confirmOrder, cancelOrder)
}

func (och OrderCommandHander) Add(commands ...Command) error {
	for _, command := range commands {
		switch command.Type {
		case "messaging.PlaceOrder": // command handle here
			{

			}
		case "messaging.ConfirmOrder": // command handle here
			{

			}
		case "messaging.CancelOrder": // command handle here
			{

			}
		}
	}
	return nil
}

func (och OrderCommandHander) Handle(command Command) {
	switch command.Type {
	case "messaging.PlaceOrder":
		och.handlePlaceOrder(command.Body.(PlaceOrder))
	case "messaging.ConfirmOrder":
		och.handleConfirmOrder(command.Body.(ConfirmOrder))
	case "messaging.CancelOrder":
		och.handleCancelOrder(command.Body.(CancelOrder))
	}
}

func (och OrderCommandHander) handlePlaceOrder(command PlaceOrder) {

}

func (och OrderCommandHander) handleConfirmOrder(command ConfirmOrder) {

}

func (och OrderCommandHander) handleCancelOrder(command CancelOrder) {

}
