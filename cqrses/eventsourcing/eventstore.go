package eventsourcing

import (
	"amaterasu/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	IEventStore interface {

		// Tries to retrieve the event sourced entity.
		// @param The id of the entity
		// @return The hydrated entity, or nil if it does not exist.
		Find(string) IEventSourced

		// Retrieves the event sourced entity.
		// @param The id of the entity
		// @return The hydrated entity, or nil if not found
		Get(string) IEventSourced

		// Saves the event sourced entity.
		// @param The entity
		// @return the error if connect fail
		Save(IEventSourced) error
	}

	// This is a basic implementation of the event store that could be optimized in the future
	// todo supports caching of snapshot in future
	EventStore struct {
		db         *gorm.DB     // using orm to easy connect
		entityType reflect.Type // entity type
		// Note: Could potentially use DataAnnotations to get a friendly/unique name in case of collisions between BCs, instead of the type's name.
		stream string
	}
)

// todo split to another class
func MysqlConnPool() *gorm.DB {

	// create connection string to mysql
	// todo make class handle many database
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/book_db"))

	// just panic error
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// config connection pool
	mysqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	mysqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	mysqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	mysqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

// Create an orm to handle save event to wvent store
// @param db connection
// @param A entity (aggreate) use this orm
// @return the orm repository
func New_EventStore(db *gorm.DB, agg IEventSourced) *EventStore {
	store := &EventStore{db: db}
	store.mapping(agg)
	return store
}

// create event store table
func (store *EventStore) CreateEventStoreTable() {
	fmt.Println(store.db.AutoMigrate(&EventData{}))
}

// golang doesn't support generic type until `golang ver 1.18` at `01/01/2022`
// map entity to ease handler
func (store *EventStore) mapping(agg IEventSourced) {
	store.entityType, store.stream = utils.GetTypeName(agg)
}

func (store *EventStore) Save(entity IEventSourced) error {
	// fmt.Println(orm.db.AutoMigrate(&EventData{}))
	events := make([]EventData, len(entity.Events()))

	for i, e := range entity.Events() {

		eTypeName := utils.GetTypeName2(reflect.TypeOf(e))
		ePayloadByte, _ := json.Marshal(e)

		eventData := EventData{
			SourceId:  e.GetSourceID(),
			Version:   e.GetVersion(),
			Type:      eTypeName,
			Payload:   ePayloadByte,
			Topic:     entity.GetTopic(),
			Partition: entity.GetPartition(),
			Stream:    store.stream,
		}
		events[i] = eventData
	}

	if errInsert := store.db.CreateInBatches(events, len(events)).Error; errInsert != nil {
		return errInsert
	}

	return nil
}

func (store *EventStore) Find(id string) IEventSourced {
	eventDataArr := []EventData{}
	store.db.
		Where(&EventData{SourceId: id, Stream: store.stream}).
		Order("version").
		Find(&eventDataArr)

	if len(eventDataArr) == 0 {
		return nil
	}

	stream := eventDataArr[0].Stream
	topic := eventDataArr[0].Topic
	partition := eventDataArr[0].Partition

	if store.stream == stream {
		entity := reflect.New(store.entityType).Interface().(IEventSourced)
		entity.CreateDefaultValue(id, topic, partition, entity)
		pastEvents := deserialize(entity, eventDataArr)
		entity.LoadFrom(pastEvents)

		return entity
	}

	return nil
}

func (orm *EventStore) Get(id string) IEventSourced {
	entity := orm.Find(id)

	if entity == nil {
		return nil
	}
	return entity
}

func deserialize(eventSourced IEventSourced, eventDataArr []EventData) []IVersionedEvent {
	events := make([]IVersionedEvent, len(eventDataArr))
	for i, eventData := range eventDataArr {
		event := reflect.New(eventSourced.GetType(eventData.Type)).Interface().(IVersionedEvent)
		json.Unmarshal(eventData.Payload, &event)
		events[i] = event
	}
	return events
}
