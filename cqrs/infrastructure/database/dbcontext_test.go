package database

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/utils"
	"testing"
)

var (
	orm = utils.MysqlConnPool("root:root@tcp(127.0.0.1:3306)/test_db")
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
	event.SetSourceId(test.Id)
	test.Events = append(test.Events, event)
}

// open this comment if you want to enable change default table
// func (InfoAggregateRoot) TableName() string {
// 	return "info_aggregate_roots" // change table name
// }

func Test_save_aggreate(t *testing.T) {
	// create event bus
	serializer := serialization.New_JsonSerializer()
	// create datacontex to save aggreate
	testDbContext := New_DataContext(orm, &InfoAggregateRoot2{}, "ni61pj1b--topic-A", 5, serializer)
	testDbContext.orm.AutoMigrate(InfoAggregateRoot2{})  // for testing purpose
	testDbContext.orm.AutoMigrate(UndispatchedMessage{}) // for testing purpose
	// create aggreate root
	agg := &InfoAggregateRoot2{}
	agg.Id = utils.NewUuidString()

	// save to database
	testDbContext.Save(agg)
}

func Test_find_aggreate_then_update_data(t *testing.T) {
	serializer := serialization.New_JsonSerializer()
	// create datacontex to save aggreate
	testDbContext := New_DataContext(orm, &InfoAggregateRoot{}, "ni61pj1b--topic-A", 5, serializer)

	testDbContext.orm.AutoMigrate(InfoAggregateRoot{}) // for testing purpose

	agg := testDbContext.Find("type id here").(*InfoAggregateRoot)

	agg.Name = "Tien"
	agg.Sex = "male"
	agg.Status = "single"

	testDbContext.Save(agg)
}

func Test_save_aggreate_then_send_events(t *testing.T) {

	// create event bus
	serializer := serialization.New_JsonSerializer()
	// create datacontex to save aggreate
	testDbContext := New_DataContext(orm, &InfoAggregateRoot{}, "ni61pj1b--topic-A", 5, serializer)
	testDbContext.orm.AutoMigrate(InfoAggregateRoot{})   // for testing purpose
	testDbContext.orm.AutoMigrate(UndispatchedMessage{}) // for testing purpose

	// create aggreate root
	agg := &InfoAggregateRoot{}
	agg.Id = utils.NewUuidString()

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
