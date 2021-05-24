package messagebroker

type Scram struct {
	Al256              bool
	Username, Password string
}

type KafkaConfig struct {
	Scr     *Scram
	Brokers string
	Topic   string
}
