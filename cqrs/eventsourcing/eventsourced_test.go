package eventsourcing

import (
	"fmt"
	"testing"
)

type Order struct {
	EventSourced
	isConfirmed  bool
	conferenceId string
	seat         int
}

func (o *Order) OnOrderPlaced(e OrderPlaced) error {
	o.isConfirmed = true
	fmt.Println(e)
	return nil
}

func (o *Order) OnOrderUpdated(e OrderUpdated) error {
	o.seat = e.Seats + 1
	fmt.Println(e)
	return nil
}

type OrderPlaced struct {
	VersionedEvent
	ConferenceId string
	AccessCode   string
	Seats        []int
}

type OrderUpdated struct {
	VersionedEvent
	Seats int
}

func Test_sth(t *testing.T) {
	o := &Order{
		EventSourced: New_EventSourced(2, "order/events", "order"),
	}
	o.AutoMappingHandles(o)

	o.Update(&OrderPlaced{AccessCode: "1234"})
}

func Test_db(t *testing.T) {
	o := &Order{
		EventSourced: New_EventSourced(1, "order/events", "order"),
	}

	o.AutoMappingHandles(o)
	o.Update(&OrderPlaced{AccessCode: "1234"})
	o.Update(&OrderUpdated{Seats: 12})

	context := New_EventSourcedORM()
	context.Create()
	context.Save(o, "1")
}
