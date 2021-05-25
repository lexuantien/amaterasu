package messagebroker

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Server struct {
	config   KafkaConfig
	consumer *kafka.Consumer
}

// :Create
func New_KafkaServer(config KafkaConfig, groupID string) (*Server, error) {

	sc := &Server{}
	sc.config = config

	config.ConfigMap.SetKey("metadata.broker.list", config.Brokers)
	config.ConfigMap.SetKey("group.id", groupID)
	config.ConfigMap.SetKey("enable.auto.commit", false)
	config.ConfigMap.SetKey("go.application.rebalance.enable", true)
	config.ConfigMap.SetKey("default.topic.config", kafka.ConfigMap{"auto.offset.reset": "earliest"})

	// just make it easy config in UI
	// can config with string external code
	if config.Scr != nil {

		var al string = "SCRAM-SHA-256"
		if !config.Scr.Al256 {
			al = "SCRAM-SHA-512"
		}

		config.ConfigMap.SetKey("security.protocol", "SASL_SSL")
		config.ConfigMap.SetKey("sasl.username", config.Scr.Username)
		config.ConfigMap.SetKey("sasl.password", config.Scr.Password)
		config.ConfigMap.SetKey("sasl.mechanisms", al)

	}

	c, err := kafka.NewConsumer(config.ConfigMap)

	if err != nil {
		panic(err)
	}

	sc.consumer = c
	err = sc.consumer.Subscribe(sc.config.Topic, nil)

	if err != nil {
		return nil, errors.New("kafka login fail")
	}

	return sc, nil
}

func (sc *Server) Receive(ctx context.Context) kafka.Event {
	// ev := <-sc.consumer.Events()
	// return ev
	return sc.consumer.Poll(100)
}

func (sc *Server) Complete(ctx context.Context, msg *kafka.Message) error {
	_, err := sc.consumer.CommitMessage(msg)
	return err
}

func (sc *Server) Assign(e kafka.AssignedPartitions) {
	sc.consumer.Assign(e.Partitions)
}

func (sc *Server) Unassign() {
	sc.consumer.Unassign()
}
