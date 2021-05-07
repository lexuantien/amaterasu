package messagebroker

import (
	"context"
	"reflect"
	"time"

	"github.com/segmentio/kafka-go"
)

// kc.Writer = kafka.NewWriter(kafka.WriterConfig{
// 	Brokers:  brokers,
// 	Balancer: &kafka.Hash{},
// 	Dialer:   dialer,
// })

type TopicClient struct {
	topic  string
	writer interface{}
}

func (c *TopicClient) SetTopic(topic string) {
	c.topic = topic
}

func (c *TopicClient) SetWriter(writer interface{}) {
	t := reflect.TypeOf(writer)
	switch t.String() {
	case "*kafka.Writer":
		c.writer = writer.(*kafka.Writer)
	}
}

func (c *TopicClient) Send(ctx context.Context, message []byte, time2live *time.Time) (bool, error) {
	var err error
	switch reflect.TypeOf(c.writer).String() {
	case "*kafka.Writer":
		if time2live != nil {
			err = c.writer.(*kafka.Writer).WriteMessages(ctx, kafka.Message{
				Topic: c.topic,
				Value: message,
				Time:  *time2live,
			})
		} else {
			err = c.writer.(*kafka.Writer).WriteMessages(ctx, kafka.Message{
				Topic: c.topic,
				Value: message,
			})
		}

	default: // hahahahha
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
