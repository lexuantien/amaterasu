package messagebroker

import (
	"context"
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type SubscriptionClient struct {
	cl     KafkaConfig
	reader *kafka.Reader
}

// :Create
func New_SubscriptionClient(cl KafkaConfig, groupID string) (*SubscriptionClient, error) {

	sc := &SubscriptionClient{}
	sc.cl = cl

	if cl.Scr != nil {
		var al scram.Algorithm = scram.SHA512
		if cl.Scr.Al256 {
			al = scram.SHA256
		}
		mechanism, err := scram.Mechanism(al, cl.Scr.Username, cl.Scr.Password)
		if err != nil {
			panic(err)
		}

		sc.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cl.Brokers,
			Topic:    cl.Topic,
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
	} else { // define later
		return nil, nil
	}

	return sc, nil
}

func (sc *SubscriptionClient) Receive(ctx context.Context) (kafka.Message, error) {
	return sc.reader.FetchMessage(ctx)
}

func (sc *SubscriptionClient) Complete(ctx context.Context, msg kafka.Message) error {
	return sc.reader.CommitMessages(ctx, msg)
}
