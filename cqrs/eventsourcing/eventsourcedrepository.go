package eventsourcing

import (
	"amaterasu/cqrs/infrastructure/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	IEventSourcedRepository interface {
		Find(string)
		Get(string)
		Save(IEventSourced, string)
	}

	EventSourcedORM struct {
		db *gorm.DB
	}
)

func New_EventSourcedORM() *EventSourcedORM {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/book_db"))

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	mysqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	mysqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	mysqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	mysqlDB.SetConnMaxLifetime(time.Hour)

	return &EventSourcedORM{
		db: db,
	}
}

func (orm EventSourcedORM) Create() {
	fmt.Println(orm.db.AutoMigrate(&EventData{}))
}

func (orm EventSourcedORM) Save(es IEventSourced, correlationId string) {

	// fmt.Println(orm.db.AutoMigrate(&EventData{}))

	events := []EventData{}

	for _, e := range es.Events() {

		eTypeName := utils.GetTypeName2(reflect.TypeOf(e))
		ePayloadByte, _ := json.Marshal(e)

		eventData := EventData{
			SourceId:  e.GetSourceID(),
			Version:   e.GetVersion(),
			Type:      eTypeName,
			Payload:   ePayloadByte,
			Topic:     es.GetTopic(),
			Partition: es.GetPartition(),
			Stream:    es.GetStream(),
			Status:    0,
		}
		events = append(events, eventData)
	}
	if errInsert := orm.db.CreateInBatches(events, len(events)).Error; errInsert != nil {
		fmt.Println(errInsert)
	}

}
