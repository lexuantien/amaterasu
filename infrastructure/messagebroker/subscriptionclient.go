package messagebroker

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SubscriptionClient struct {
	cl       KafkaConfig
	consumer *kafka.Consumer
}

// :Create
func New_SubscriptionClient(cl KafkaConfig, groupID string) (*SubscriptionClient, error) {

	sc := &SubscriptionClient{}
	sc.cl = cl

	config := &kafka.ConfigMap{}

	if cl.Scr != nil {
		var al string = "SCRAM-SHA-256"
		if !cl.Scr.Al256 {
			al = "SCRAM-SHA-512"
		}
		config = &kafka.ConfigMap{
			"metadata.broker.list":            cl.Brokers,
			"security.protocol":               "SASL_SSL",
			"sasl.mechanisms":                 al,
			"sasl.username":                   cl.Scr.Username,
			"sasl.password":                   cl.Scr.Password,
			"group.id":                        groupID,
			"go.events.channel.enable":        true,
			"go.application.rebalance.enable": true,
			// "enable.auto.commit":              false,
			"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
			// "debug":                           "generic,broker,security",
		}
	} else { // define later
		return nil, errors.New("only support kafka scram connect")
	}

	c, err := kafka.NewConsumer(config)

	if err != nil {
		panic(err)
	}

	sc.consumer = c
	err = sc.consumer.Subscribe(sc.cl.Topic, nil)

	if err != nil {
		return nil, errors.New("kafka login fail")
	}

	return sc, nil
}

func (sc *SubscriptionClient) Receive(ctx context.Context) kafka.Event {
	ev := <-sc.consumer.Events()
	return ev
	// return sc.consumer.Poll(100)
}

func (sc *SubscriptionClient) Complete(ctx context.Context, msg *kafka.Message) error {
	_, err := sc.consumer.CommitMessage(msg)
	return err
}

func (sc *SubscriptionClient) Assign(e kafka.AssignedPartitions) {
	sc.consumer.Assign(e.Partitions)
}

func (sc *SubscriptionClient) Unassign() {
	sc.consumer.Unassign()
}
