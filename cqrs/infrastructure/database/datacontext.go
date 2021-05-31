package database

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/utils"
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

func New_DataContext(table string, agg IAggregateRoot, bus messaging.IEventBus) *DataContext {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/" + table))
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
	agg := reflect.New(context.aggType).Interface().(IAggregateRoot)
	err := context.db.Where("id = ?", id).First(&agg).Error
	if err != nil {
		panic(err)
	}
	return agg
}

func (context *DataContext) Save(aggreate IAggregateRoot) {
	// aggreate.GetID()

	agg := reflect.New(context.aggType).Interface().(IAggregateRoot)
	context.db.Where(aggreate).First(&agg)

	if agg != nil {
		return
	}
	// Can't have transactions across storage and message bus.
	err := context.db.Save(aggreate).Error
	if err != nil {
		panic(err)
	}

	if val, ok := reflect.TypeOf(aggreate).(messaging.IEventPublisher); ok {
		envelopes := make([]messaging.Envelope, len(val.GetEvents()))

		// all event need 1 partition
		// random choose partition key
		partititonKey := rand.Intn(5)
		for i, e := range val.GetEvents() {
			envelopes[i] = messaging.EnvelopeWrap2(e, int32(partititonKey), messaging.EVENT)
		}
		context.bus.Publishes(envelopes...)
	}
}
