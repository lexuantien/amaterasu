package database

type MessageAction uint

const (
	EVENT MessageAction = iota
	COMMAND
)

type UndispatchedMessage struct {
	Id           string        `gorm:"column:source_id"`
	Stream       string        `gorm:"column:stream"`
	Body         []byte        `gorm:"column:payload;type:JSON"`
	PartitionKey int32         `gorm:"column:key"`
	Topic        string        `gorm:"column:topic"`
	MsgType      string        `gorm:"column:msg_type"`
	MsgAction    MessageAction `gorm:"column:msg_action"`
}

func (UndispatchedMessage) TableName() string {
	return "undispatched_message"
}
