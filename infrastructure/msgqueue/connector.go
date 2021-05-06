package msgqueue

type Algorithm int
type Queue int

const (
	SHA216   Algorithm = 0
	SHA512   Algorithm = 1
	KAFKA    Queue     = 2
	RABBITMQ Queue     = 3
)

type Connector struct {
	scr     Scram
	brokers []string
	//
	cnn interface{}
}

type Scram struct {
	Username string
	Password string
	_256     bool
}

func (c *Connector) Connect(scr Scram, brokers []string, queue Queue) {
	c.scr = scr
	c.brokers = brokers

	switch queue {
	case KAFKA:
		c.cnn.(Kafka).connect(scr, brokers)
	default: // hahahahha
	}
}
