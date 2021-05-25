package handling

import (
	"fmt"
	kafkaa "leech-service/cqrs/infrastructure/kafkaa"
	"sync"
)

//! config
var (
	wg          = sync.WaitGroup{}
	kafkaConfig = kafkaa.KafkaConfig{
		Scr: &kafkaa.Scram{
			Username: "ni61pj1b",
			Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
			Al256:    true,
		},
		Brokers: "glider-01.srvs.cloudkafka.com:9094,glider-02.srvs.cloudkafka.com:9094,glider-03.srvs.cloudkafka.com:9094",
		Topic:   "ni61pj1b--topic-A",
	}
)

//! command handlers
type CommandHander1 struct{}
type CommandHander2 struct{}

func New_CommandHander1() *CommandHander1 {
	return &CommandHander1{}
}

func New_CommandHander2() *CommandHander2 {
	return &CommandHander2{}
}

//* command handler func
func (o CommandHander1) Handle(command interface{}) error {
	switch cmd := command.(type) {

	case *Foo1:
		fmt.Println(cmd)
	case *Foo2:
	case *Foo3:

	default:
		panic("Command not found in this command handler")
	}

	return nil
}

func (o CommandHander2) Handle(command interface{}) error {
	switch command.(type) {

	case *Foo1:
		fmt.Println(command.(*Foo1))

	default:
		panic("Command not found in this command handler")
	}

	return nil
}

//? commands
type Foo1 struct {
	ProductId   string
	Quantity    uint
	Description string
}

type Foo2 struct {
	OrderId string
}

type Foo3 struct {
	OrderId string
}

type Foo4 struct{}

//? events
type Event1 struct {
	name string
}

type Event3 struct {
	name string
	old  int
	n    string
}

type Event2 struct {
	name string
	div  int
}

//! command handlers
type EventHandler1 struct{}
type EventHandler2 struct{}

func New_EventHandler1() *EventHandler1 {
	return &EventHandler1{}
}
func New_EventHandler2() *EventHandler2 {
	return &EventHandler2{}
}

func (e EventHandler1) Handle(ev interface{}) error {

	return nil
}

func (e EventHandler2) Handle(ev interface{}) error {

	return nil
}
