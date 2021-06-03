package database

type (
	IAggregateRoot interface {
		GetId() string
		SetId(string)
		GetVersion() int64
		SetVersion(int64)
	}

	AggregateRoot struct {
		Id      string `gorm:"column:id"`
		Version int64  `gorm:"column:version"`
	}
)

func (agg *AggregateRoot) GetId() string {
	return agg.Id
}

func (agg *AggregateRoot) SetId(id string) {
	agg.Id = id
}

func (agg *AggregateRoot) GetVersion() int64 {
	return agg.Version
}

func (agg *AggregateRoot) SetVersion(newVersion int64) {
	agg.Version = newVersion
}
