package messagebroker

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type TopicClient struct {
	cl     KafkaConfig
	writer *kafka.Writer
}

func New_TopicClient(cl KafkaConfig) *TopicClient {

	tc := &TopicClient{}
	tc.cl = cl

	if cl.Scr != nil {
		var al scram.Algorithm = scram.SHA512
		if cl.Scr.Al256 {
			al = scram.SHA256
		}
		mechanism, err := scram.Mechanism(al, cl.Scr.Username, cl.Scr.Password)
		if err != nil {
			panic(err)
		}

		dialer := &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
			SASLMechanism: mechanism,
		}

		tc.writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers:      cl.Brokers,
			Balancer:     &kafka.Hash{},
			Dialer:       dialer,
			RequiredAcks: int(kafka.RequireOne),
		})
	} else {

	}

	return tc
}

func (c *TopicClient) Send(ctx context.Context, message, clazzType []byte) error {
	return c.writer.WriteMessages(ctx, kafka.Message{
		Topic: c.cl.Topic,
		Key:   clazzType,
		Value: message,
	})
}
