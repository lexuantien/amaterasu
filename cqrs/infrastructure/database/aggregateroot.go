package database

type (
	IAggregateRoot interface {
		GetID() string
	}

	AggregateRoot struct {
		ID string
	}
)

func (agg *AggregateRoot) GetID() string {
	return agg.ID
}
