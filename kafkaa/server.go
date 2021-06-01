package kafkaa

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// kafka consumer
type Server struct {
	config   KafkaConfig
	consumer *kafka.Consumer
}

// create kafka consumer configuration
// config kafka config
// 	groupID consumer group id
// 	Server consumer config
// 	error nil if success, else
/// exclude metadata.broker.list, group.id, go.application.rebalance.enable, default.topic.config
func New_ConsumerConfig(config KafkaConfig, groupID string) (*Server, error) {

	sc := &Server{}
	config.configMap = &kafka.ConfigMap{}
	sc.config = config

	config.configMap.SetKey("metadata.broker.list", config.Brokers)
	config.configMap.SetKey("group.id", groupID)
	config.configMap.SetKey("enable.auto.commit", false)
	config.configMap.SetKey("go.application.rebalance.enable", true)
	config.configMap.SetKey("default.topic.config", kafka.ConfigMap{"auto.offset.reset": "earliest"})

	for k, v := range config.ConfigMap {
		config.configMap.SetKey(k, v)
	}

	// start consume kafka
	consumer, err := kafka.NewConsumer(config.configMap)

	if err != nil {
		return nil, err
	}

	sc.consumer = consumer
	err = sc.consumer.Subscribe(sc.config.Topic, nil)

	if err != nil {
		return nil, errors.New("create kafka subscribe fail")
	}

	return sc, nil
}

// receive message
// kafka.Event: event store kafka message
func (sc *Server) Receive(ctx context.Context) kafka.Event {
	// ev := <-sc.consumer.Events()
	// return ev
	return sc.consumer.Poll(10) // poll message each 10 ms
}

// commit message
// error nil or error
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
