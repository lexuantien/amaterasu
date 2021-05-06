package msgqueue

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Kafka struct {
	Writer *kafka.Writer
}

func (kc Kafka) connect(scr Scram, brokers []string) {

	//
	var mechanism sasl.Mechanism
	var err error

	if scr._256 {
		mechanism, err = scram.Mechanism(scram.SHA256, scr.Username, scr.Password)
	} else {
		mechanism, err = scram.Mechanism(scram.SHA512, scr.Username, scr.Password)

	}

	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	kc.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})
}
