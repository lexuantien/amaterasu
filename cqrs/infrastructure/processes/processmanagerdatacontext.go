package processes

import (
	"amaterasu/cqrs/infrastructure/database"
	"amaterasu/cqrs/infrastructure/serialization"
	"amaterasu/utils"
	"math/rand"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type (
	IProcessManagerDataContext interface {
		Find(string, bool) IProcessManager
		Save(IProcessManager) error
	}

	ProcessManagerDataContext struct {
		topic         string
		numPartittion int
		serializer    serialization.ISerializer
		pmType        reflect.Type
		orm           *gorm.DB
	}
)

func New_ProcessManagerDataContext(db *gorm.DB, pm interface{}, topic string, numPartittion int, serializer serialization.ISerializer) *ProcessManagerDataContext {
	// todo split to another class

	return &ProcessManagerDataContext{
		orm:           db,
		topic:         topic,
		numPartittion: numPartittion,
		serializer:    serializer,
		pmType:        reflect.TypeOf(pm),
	}
}

func (pmContext *ProcessManagerDataContext) Find(id string, completed bool) IProcessManager {
	pm := reflect.New(pmContext.pmType.Elem()).Interface().(IProcessManager)

	tx := pmContext.orm.Model(pm).Where("id = ? and completed = ?", id, completed).First(pm)
	if tx.RowsAffected != 1 {
		return nil
	}

	return pm
}

func (pmContext *ProcessManagerDataContext) Save(pm IProcessManager) error {

	oldVersion := pm.GetVersion()
	pm.SetVersion(time.Now().Unix())

	// begin a transaction
	tx := pmContext.orm.Begin()

	if err := tx.Save(pm).Where("version = ?", oldVersion).Error; err != nil {
		tx.Rollback()
		return err
	}

	commands := pm.GetCommands()
	if len(commands) > 0 {
		undispatchedArr := make([]database.UndispatchedMessage, len(commands))

		// all command need 1 partition
		// random choose partition key
		partititonKey := rand.Intn(pmContext.numPartittion)
		for i, cmd := range commands {
			bodyByte, err := pmContext.serializer.Serialize(cmd)
			if err != nil {
				tx.Rollback()
			}
			undispatchedArr[i] = database.UndispatchedMessage{
				SourceId:  pm.GetId(),
				Stream:    utils.GetObjType2(pmContext.pmType.Elem()),
				Payload:   bodyByte,
				Partition: int32(partititonKey),
				Topic:     pmContext.topic,
				MsgAction: database.COMMAND,
				MsgType:   utils.GetObjType2(reflect.TypeOf(cmd)),
				Status:    0,
			}
		}
		if errInsert := tx.CreateInBatches(undispatchedArr, len(undispatchedArr)).Error; errInsert != nil {
			tx.Rollback()
			return errInsert
		}
	}

	// end a transaction
	return tx.Commit().Error
}
