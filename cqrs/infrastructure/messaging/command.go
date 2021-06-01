package messaging

type (
	ICommandHandler interface{}

	ICommandHandlerRegistry interface {
		Register(ICommandHandler) error
	}

	ICommand interface {
		GetId() string
		SetId(string)
	}

	Command struct {
		Id string
	}
)

func (cmd *Command) GetId() string {
	return cmd.Id
}

func (cmd *Command) SetId(id string) {
	cmd.Id = id
}
