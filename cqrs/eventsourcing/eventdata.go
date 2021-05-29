package eventsourcing

type EventData struct {
	SourceId  string `gorm:"column:source_id"`
	Stream    string `gorm:"column:stream"`
	Topic     string `gorm:"column:topic"`
	Partition int    `gorm:"column:partition"`
	Version   int    `gorm:"column:version"`
	Type      string `gorm:"column:type"`
	Status    int    `gorm:"column:status"`
	Payload   []byte `gorm:"column:payload"`
	Error     string `gorm:"column:error"`
}

func (EventData) TableName() string {
	return "event_store"
}
