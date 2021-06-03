package eventsourcing

import (
	"amaterasu/utils"
	"fmt"
	"testing"
)

var mysqlPool = utils.MysqlConnPool("root:root@tcp(127.0.0.1:3306)/test_db")

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
		EventSourced: New_EventSourced(2, "order/events"),
	}
	o.AutoMappingHandles(o)

	o.Update(&OrderPlaced{AccessCode: "1234"})
}

func Test_save(t *testing.T) {
	o := &Order{
		EventSourced: New_EventSourced(1, "order/events"),
	}

	o.AutoMappingHandles(o)
	o.Update(&OrderPlaced{AccessCode: "1234"})
	o.Update(&OrderUpdated{Seats: 12})
	orm := New_EventStore(mysqlPool, &Order{})
	orm.CreateEventStoreTable()
	orm.Save(o)
	// orm.Find("d6003786-5b12-46d0-ba3a-5b5fadcde339")
}

func Test_find(t *testing.T) {
	o := &Order{
		EventSourced: New_EventSourced(1, "order/events"),
	}
	o.AutoMappingHandles(o)
	// o.Update(&OrderPlaced{AccessCode: "1234"})
	// o.Update(&OrderUpdated{Seats: 12})
	orm := New_EventStore(mysqlPool, &Order{})
	agg := orm.Find("9a7815ba-d4f0-4c86-b564-1624871bd24e")
	fmt.Println(agg)
}
