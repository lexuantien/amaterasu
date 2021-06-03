package processes

import "amaterasu/cqrs/infrastructure/messaging"

type (
	IProcessManager interface {
		GetId() string
		IsComplete() bool
		GetCommands() []messaging.ICommand
		GetVersion() int64
		SetVersion(int64)
	}

	ProcessManager struct {
		Id        string               `gorm:"column:id"`
		Version   int64                `gorm:"column:version"` // optimistic lock, using time.Now().Unix()
		Completed bool                 `gorm:"column:completed"`
		commands  []messaging.ICommand `gorm:"-"`
	}
)

func (pm *ProcessManager) GetId() string {
	return pm.Id
}

func (pm *ProcessManager) IsComplete() bool {
	return pm.Completed
}

func (pm *ProcessManager) GetCommands() []messaging.ICommand {
	return pm.commands
}

func (pm *ProcessManager) GetVersion() int64 {
	return pm.Version
}

func (pm *ProcessManager) SetVersion(newVersion int64) {
	pm.Version = newVersion
}
