package handling

import "fmt"

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
