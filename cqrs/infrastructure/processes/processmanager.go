package processes

import "amaterasu/cqrs/infrastructure/messaging"

type (
	IProcessManager interface {
		GetId() string

		IsComplete() bool

		GetCommands() []messaging.ICommand
	}
)
