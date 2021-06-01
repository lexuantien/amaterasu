package processes

type (
	IProcessManagerDataContext interface {
		Find(id string)
		Save(IProcessManager)
	}
)
