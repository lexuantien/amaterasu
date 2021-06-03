package database

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/utils"
	"math/rand"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type (

	// interface contains func handle store/load aggreate in cqrs
	IDataContext interface {
		// find an aggreate in database
		// id the aggreate id
		// aggreate or nil if not found
		Find(id string) IAggregateRoot

		// save aggreate to database
		// aggreate the aggreate
		// nil or error if not found
		Save(aggreate IAggregateRoot) error
	}

	// handle store/load aggreate in cqrs
	DataContext struct {
		// eventbus to publish message
		// bus     messaging.IEventBus
		orm           *gorm.DB
		aggType       reflect.Type
		topic         string
		numPartittion int
		serializer    serialization.ISerializer
	}
)

func New_DataContext(db *gorm.DB, agg interface{}, topic string, numPartittion int, serializer serialization.ISerializer) *DataContext {
	// todo split to another class

	return &DataContext{
		orm:           db,
		topic:         topic,
		numPartittion: numPartittion,
		serializer:    serializer,
		aggType:       reflect.TypeOf(agg),
	}
}

func (context *DataContext) Find(id string) IAggregateRoot {

	agg := reflect.New(context.aggType.Elem()).Interface().(IAggregateRoot)

	tx := context.orm.Model(agg).Where("id=?", id).First(agg)
	if tx.Error != nil || tx.RowsAffected != 1 {
		return nil
	}

	return agg
}

func (dbContext *DataContext) Save(aggreate IAggregateRoot) error {

	oldVersion := aggreate.GetVersion()
	aggreate.SetVersion(time.Now().Unix())

	// begin a transaction
	tx := dbContext.orm.Begin()

	// Can't have transactions across storage and UndispatchedMessage.
	if err := tx.Save(aggreate).Where("version = ?", oldVersion).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Can't have transactions across storage and UndispatchedMessage.
	if val, ok := aggreate.(messaging.IEventPublisher); ok && val.GetEvents() != nil {
		undispatchedArr := make([]UndispatchedMessage, len(val.GetEvents()))

		// all event need 1 partition
		// random choose partition key
		partititonKey := rand.Intn(dbContext.numPartittion)
		for i, e := range val.GetEvents() {
			bodyByte, err := dbContext.serializer.Serialize(e)
			if err != nil {
				tx.Rollback()
				return err
			}
			undispatchedArr[i] = UndispatchedMessage{
				SourceId:  aggreate.GetId(),
				Stream:    utils.GetObjType2(dbContext.aggType.Elem()),
				Payload:   bodyByte,
				Partition: int32(partititonKey),
				Topic:     dbContext.topic,
				MsgAction: EVENT,
				MsgType:   utils.GetObjType2(reflect.TypeOf(e)),
				Status:    0,
			}
		}

		if errInsert := dbContext.orm.CreateInBatches(undispatchedArr, len(undispatchedArr)).Error; errInsert != nil {
			tx.Rollback()
			return errInsert
		}
	}

	return tx.Commit().Error
}
