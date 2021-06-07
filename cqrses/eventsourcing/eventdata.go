package eventsourcing

import (
	"time"
)

type EventData struct {
	SourceId  string `gorm:"column:source_id;primaryKey"`
	Stream    string `gorm:"column:stream;primaryKey"`
	Version   int    `gorm:"column:version;primaryKey"`
	Partition int32  `gorm:"column:partition"`
	Topic     string `gorm:"column:topic"`
	EventId   uint   `gorm:"column:event_id"`
	Type      string `gorm:"column:type"`
	Payload   []byte `gorm:"column:payload;type:JSON"`
	CreatedAt time.Time
}

func (EventData) TableName() string {
	return "event_store"
}
