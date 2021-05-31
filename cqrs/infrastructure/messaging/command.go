package messaging

type (
	ICommandHandler interface{}

	ICommandHandlerRegistry interface {
		Register(ICommandHandler) error
	}

	ICommand interface {
		GetID() string
		SetID(string)
	}

	Command struct {
		ID string
	}
)

func (c *Command) GetID() string {
	return c.ID
}

func (c *Command) SetID(id string) {
	c.ID = id
}
