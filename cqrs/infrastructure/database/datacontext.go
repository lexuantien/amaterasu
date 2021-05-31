package database

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/utils"
	"context"
	"math/rand"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	IDataContext interface {
		Find(id string)
		Save(aggreate IAggregateRoot)
	}

	DataContext struct {
		bus         messaging.IEventBus
		db          *gorm.DB
		aggType     reflect.Type
		aggTypeName string
	}
)

func New_DataContext(agg interface{}, bus messaging.IEventBus) *DataContext {
	// todo split to another class
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test_db"))
	if err != nil {
		panic(err)
	}

	return &DataContext{
		db:          db,
		bus:         bus,
		aggType:     reflect.TypeOf(agg),
		aggTypeName: utils.GetTypeName2(reflect.TypeOf(agg)),
	}
}

func (context *DataContext) Find(id string) IAggregateRoot {

	agg := reflect.New(context.aggType.Elem()).Interface().(IAggregateRoot)

	tx := context.db.Model(agg).Where("id=?", id).First(agg)
	if tx.Error != nil {
		panic(tx.Error)
	}

	return agg
}

func (dbContext *DataContext) Save(aggreate IAggregateRoot) {
	// aggreate.GetID()

	// agg := reflect.New(dbContext.aggType).Interface().(IAggregateRoot)
	// db := dbContext.db.Model(aggreate).Where(aggreate).First(&agg)
	// fmt.Println(db.RowsAffected)

	// if db.RowsAffected != 0 {
	// 	return
	// }

	// Can't have transactions across storage and message bus.
	err := dbContext.db.Save(aggreate).Error
	if err != nil {
		panic(err)
	}

	// if contains some event need to publish
	// todo store to `undispatch event` table then listener it to publish
	if val, ok := aggreate.(messaging.IEventPublisher); ok && val.GetEvents() != nil {
		envelopes := make([]messaging.Envelope, len(val.GetEvents()))
		// all event need 1 partition
		// random choose partition key
		partititonKey := rand.Intn(5)
		for i, e := range val.GetEvents() {
			envelopes[i] = messaging.EnvelopeWrap2(e, int32(partititonKey), messaging.EVENT)
		}
		dbContext.bus.Publishes(context.Background(), envelopes...)
	}
}
