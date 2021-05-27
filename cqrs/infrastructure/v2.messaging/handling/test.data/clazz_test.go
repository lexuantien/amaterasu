package handling

import (
	kafkaa "amaterasu/cqrs/infrastructure/kafkaa"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
		ConfigMap: &kafka.ConfigMap{},
	}
)

//? commands
type Foo1 struct {
	ProductId   string
	Quantity    uint
	Description string
}

func (f Foo1) Id() string {
	return f.ProductId
}

type Foo2 struct {
	OrderId string
}

func (f Foo2) Id() string {
	return f.OrderId
}

type Foo3 struct {
	FId    string
	Gender bool
}

func (f Foo3) Id() string {
	return f.FId
}

type Foo4 struct {
	FId string
	T   string
	N   bool
}

func (f Foo4) Id() string {
	return f.FId
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
	Name string
	Id   string
}

type Event2 struct {
	Id   string
	Name string
	Div  int
}

type Event3 struct {
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
