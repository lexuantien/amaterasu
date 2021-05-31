package database

type (
	IAggregateRoot interface {
		GetID() string
		SetID(string)
	}

	AggregateRoot struct {
		ID string `gorm:"column:id"`
	}
)

func (agg *AggregateRoot) GetID() string {
	return agg.ID
}

func (agg *AggregateRoot) SetID(id string) {
	agg.ID = id
}
