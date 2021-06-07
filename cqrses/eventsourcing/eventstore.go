package eventsourcing

import (
	"amaterasu/cqrs/infrastructure/database"
	"amaterasu/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"

	"gorm.io/gorm"
)

type (
	IEventStore interface {

		// Tries to retrieve the event sourced entity.
		// @param The id of the entity
		// @return The hydrated entity, or nil if it does not exist.
		Find(string) IEventSourced

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
		stream       string
		topic        string
		numPartition int32
	}
)

// Create an orm to handle save event to wvent store
// @param db connection
// @param A entity (aggreate) use this orm
// @return the orm repository
func New_EventStore(db *gorm.DB, agg IEventSourced, numPartition int32, topic string) *EventStore {
	store := &EventStore{
		db:           db,
		numPartition: numPartition,
		topic:        topic,
	}
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
	store.entityType, store.stream = utils.GetObjType(agg)
}

func (store *EventStore) Save(entity IEventSourced) error {
	// fmt.Println(orm.db.AutoMigrate(&EventData{}))
	events := make([]EventData, len(entity.Events()))
	undispatchMessages := make([]database.UndispatchedMessage, len(entity.Events()))

	var partition int32 = int32(rand.Intn(int(store.numPartition)))
	if entity.GetPartition() != -1 {
		partition = entity.GetPartition()
	}

	tx := store.db.Begin()
	for i, e := range entity.Events() {

		eTypeName := utils.GetObjType2(reflect.TypeOf(e))
		ePayloadByte, err := json.Marshal(e)
		if err != nil {
			return err
		}
		// `event_store` table
		eventData := EventData{
			SourceId:  e.GetSourceId(),
			Stream:    store.stream,
			Version:   e.GetVersion(),
			Topic:     store.topic,
			Partition: partition,
			Type:      eTypeName,
			Payload:   ePayloadByte,
		}

		// `undispatched_message` table
		undispatchedMessage := database.UndispatchedMessage{
			SourceId:  e.GetSourceId(),
			Stream:    store.stream,
			Partition: partition,
			Topic:     store.topic,
			Payload:   ePayloadByte,
			MsgType:   eTypeName,
			MsgAction: database.EVENT,
			Status:    0,
			Version:   entity.GetVersion(),
		}

		events[i] = eventData
		undispatchMessages[i] = undispatchedMessage
	}

	if err := tx.CreateInBatches(events, len(events)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if errInsert := tx.CreateInBatches(events, len(events)).Error; errInsert != nil {
		tx.Rollback()
		return errInsert
	}

	return tx.Commit().Error
}

func (store *EventStore) Find(id string) IEventSourced {
	var eventDataArr []EventData
	result := store.db.
		Where(&EventData{SourceId: id, Stream: store.stream}).
		Order("version").
		Find(&eventDataArr)

	if result.RowsAffected == 0 {
		return nil
	}

	entity := reflect.New(store.entityType).Interface().(IEventSourced)
	entity.CreateDefaultValue(id, store.topic, eventDataArr[0].Partition, entity)
	// entity.SetPartition(eventDataArr[0].Partition)
	pastEvents := deserialize(entity, eventDataArr)
	entity.LoadFromHistory(pastEvents)

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
