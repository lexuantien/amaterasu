package database

type (
	IAggregateRoot interface {
		GetId() string
		SetId(string)
	}

	AggregateRoot struct {
		Id string `gorm:"column:id"`
	}
)

func (agg *AggregateRoot) GetId() string {
	return agg.Id
}

func (agg *AggregateRoot) SetId(id string) {
	agg.Id = id
}
