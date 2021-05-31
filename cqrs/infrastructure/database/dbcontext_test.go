package database

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/kafkaa"
	"amaterasu/utils"
	"testing"
)

var (
	kafkaConfig = kafkaa.KafkaConfig{
		Scr: &kafkaa.Scram{
			Username: "ni61pj1b",
			Password: "mWl_TWtiOPUKF4hRXVXPfULsKSoMzT0l",
			Al256:    true,
		},
		Brokers:   "glider-01.srvs.cloudkafka.com:9094,glider-02.srvs.cloudkafka.com:9094,glider-03.srvs.cloudkafka.com:9094",
		Topic:     "ni61pj1b--topic-A",
		ConfigMap: make(map[string]interface{}),
	}
)

type (
	NameChanged struct {
		messaging.Event
	}

	StatusChanged struct {
		messaging.Event
	}

	SexChanged struct {
		messaging.Event
	}

	InfoAggregateRoot struct {
		messaging.EventPublisher // publish event
		AggregateRoot            //

		Name   string `gorm:"column:name"`
		Status string `gorm:"column:status"`
		Sex    string `gorm:"column:sex"`
	}

	InfoAggregateRoot2 struct {
		AggregateRoot
		//
		Name   string `gorm:"column:name"`
		Status string `gorm:"column:status"`
		Sex    string `gorm:"column:sex"`
	}
)

func (test *InfoAggregateRoot) AddEvent(event messaging.IEvent) {
	event.SetSourceID(test.ID)
	test.Events = append(test.Events, event)
}

// open this comment if you want to enable change default table
// func (InfoAggregateRoot) TableName() string {
// 	return "info_aggregate_roots" // change table name
// }

func Test_save_aggreate(t *testing.T) {
	// create event bus
	kafkaConfig := kafkaa.New_kafkaa(kafkaConfig)
	sender := messaging.New_Producer(*kafkaConfig)
	serializer := serialization.New_JsonSerializer()
	bus := messaging.New_EventBus(sender, serializer)

	// create datacontex to save aggreate
	testDbContext := New_DataContext(&InfoAggregateRoot2{}, bus)
	testDbContext.db.AutoMigrate(InfoAggregateRoot2{})

	// create aggreate root
	agg := &InfoAggregateRoot2{}
	agg.ID = utils.NewString()

	// save to database
	testDbContext.Save(agg)
}

func Test_find_aggreate_then_update_data(t *testing.T) {
	// create event bus
	kafkaConfig := kafkaa.New_kafkaa(kafkaConfig)
	sender := messaging.New_Producer(*kafkaConfig)
	serializer := serialization.New_JsonSerializer()
	bus := messaging.New_EventBus(sender, serializer)

	// create datacontex to save aggreate
	testDbContext := New_DataContext(&InfoAggregateRoot{}, bus)
	testDbContext.db.AutoMigrate(InfoAggregateRoot{}) // for testing purpose

	agg := testDbContext.Find("type id here").(*InfoAggregateRoot)

	agg.Name = "Tien"
	agg.Sex = "male"
	agg.Status = "single"

	testDbContext.Save(agg)
}

func Test_save_aggreate_then_send_events(t *testing.T) {

	// create event bus
	kafkaConfig := kafkaa.New_kafkaa(kafkaConfig)
	sender := messaging.New_Producer(*kafkaConfig)
	serializer := serialization.New_JsonSerializer()
	bus := messaging.New_EventBus(sender, serializer)

	// create datacontex to save aggreate
	testDbContext := New_DataContext(&InfoAggregateRoot{}, bus)
	testDbContext.db.AutoMigrate(InfoAggregateRoot{})

	// create aggreate root
	agg := &InfoAggregateRoot{}
	agg.ID = utils.NewString()

	// create events
	e1 := &NameChanged{}
	e2 := &SexChanged{}
	e3 := &StatusChanged{}

	// add event to aggreate
	agg.AddEvent(e1)
	agg.AddEvent(e2)
	agg.AddEvent(e3)

	// save to database
	testDbContext.Save(agg)
}
