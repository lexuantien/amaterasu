package database

type MessageAction uint

const (
	EVENT MessageAction = iota
	COMMAND
)

type UndispatchedMessage struct {
	SourceId  string        `gorm:"column:source_id"`
	Stream    string        `gorm:"column:stream"`
	Version   int           `gorm:"column:version"`
	Payload   []byte        `gorm:"column:payload;type:JSON"`
	Partition int32         `gorm:"column:partition"`
	Topic     string        `gorm:"column:topic"`
	MsgType   string        `gorm:"column:msg_type"`
	MsgAction MessageAction `gorm:"column:msg_action"`
	Status    int           `gorm:"column:status;comment:'-1=FAILURE|0=NEW|1=RUNNING|2=SUCCESS'"`
	Error     string        `gorm:"column:error"`
}

func (UndispatchedMessage) TableName() string {
	return "undispatched_message"
}
