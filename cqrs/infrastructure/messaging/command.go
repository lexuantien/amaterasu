package messaging

type (
	ICommandHandler interface{}

	ICommandHandlerRegistry interface {
		Register(ICommandHandler) error
	}

	// command interface
	ICommand interface {
		GetId() string
		SetId(string)
	}

	// command struct
	Command struct {
		// the id or source id maybe
		Id string
	}
)

func (cmd *Command) GetId() string {
	return cmd.Id
}

func (cmd *Command) SetId(id string) {
	cmd.Id = id
}
