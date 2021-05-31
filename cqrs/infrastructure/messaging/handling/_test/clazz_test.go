package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/kafkaa"
	"fmt"
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
		Brokers:   "glider-01.srvs.cloudkafka.com:9094,glider-02.srvs.cloudkafka.com:9094,glider-03.srvs.cloudkafka.com:9094",
		Topic:     "ni61pj1b--topic-A",
		ConfigMap: make(map[string]interface{}),
	}
)

//? commands
type Foo1 struct {
	messaging.Command
	ProductId   string
	Quantity    uint
	Description string
}

type Foo2 struct {
	messaging.Command
	OrderId string
}

type Foo3 struct {
	messaging.Command
	FId    string
	Gender bool
}

type Foo4 struct {
	messaging.Command
	FId string
	T   string
	N   bool
}

//! command handlers
type CommandHander1 struct{}
type CommandHander2 struct{}

func New_CommandHander1() *CommandHander1 {
	return &CommandHander1{}
}

func (o CommandHander1) HandleFoo1(command Foo1) error {
	fmt.Println(command)
	return nil
}

func (o CommandHander1) HandleFoo2(command Foo2) error {
	fmt.Println(command)
	return nil
}

func (o CommandHander1) HandleFoo3(command Foo3) error {
	fmt.Println(command)
	return nil
}

func (o CommandHander2) HandleFoo4(command Foo4) error {
	fmt.Println(command)
	return nil
}

func New_CommandHander2() *CommandHander2 {
	return &CommandHander2{}
}

//? events
type Event1 struct {
	messaging.Event
	Name string
	Id   string
}

type Event2 struct {
	messaging.Event
	Id   string
	Name string
	Div  int
}

type Event3 struct {
	messaging.Event
	Id   string
	Name string
	Old  int
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

func (eh EventHandler1) HandleEvent1(e Event1) error {
	fmt.Println("EventHandler1 - ", e)
	return nil
}

func (eh EventHandler1) HandleEvent2(e Event2) error {
	fmt.Println("EventHandler1 - ", e)
	return nil
}

//

func (eh EventHandler2) HandleEvent1(e Event1) error {
	fmt.Println("EventHandler2 - ", e)
	return nil
}

func (eh EventHandler2) HandleEvent2(e Event2) error {
	fmt.Println("EventHandler2 - ", e)
	return nil
}

func (eh EventHandler2) HandleEvent3(e Event3) error {
	fmt.Println("EventHandler2 - ", e)
	return nil
}
