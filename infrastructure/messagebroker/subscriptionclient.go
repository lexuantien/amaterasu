package messagebroker

import (
	"context"
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Scram struct {
	Al256              bool
	Username, Password string
}

type SubscriptionClient struct {
	reader *kafka.Reader
}

// :Create
func New_SubscriptionClient(scr Scram, brokers []string, topic, groupID string) (*SubscriptionClient, error) {

	sc := &SubscriptionClient{}
	var al scram.Algorithm = scram.SHA512
	if scr.Al256 {
		al = scram.SHA256
	}
	mechanism, err := scram.Mechanism(al, scr.Username, scr.Password)
	if err != nil {
		panic(err)
	}

	sc.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID, // Send messages to only one subscriber per group.
		MinBytes: 10e3,    // 10KB
		MaxBytes: 10e6,    // 10MB
		Dialer: &kafka.Dialer{
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
			SASLMechanism: mechanism,
		},
	})

	return sc, nil
}

func (sc *SubscriptionClient) Receive(ctx context.Context) (kafka.Message, error) {
	return sc.reader.FetchMessage(ctx)
}

func (sc *SubscriptionClient) Complete(ctx context.Context, msg kafka.Message) error {
	return sc.reader.CommitMessages(ctx, msg)
}
