package messagebroker

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Scram struct {
	Al256              bool
	Username, Password string
}

type SubscriptionClient struct {
	reader *kafka.Reader
	wg     sync.WaitGroup
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

	joined := make(chan struct{})
	sc.wg.Add(1)
	sc.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		Topic:                  topic,
		GroupID:                groupID,     // Send messages to only one subscriber per group.
		MaxBytes:               100e3,       // 100KB
		MaxWait:                time.Second, // Allow to exit readloop in max 1s.
		PartitionWatchInterval: time.Second,
		WatchPartitionChanges:  true,
		StartOffset:            kafka.LastOffset, // Don't read old messages.
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
			SASLMechanism: mechanism,
		},
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// NOTE: Hacky way to use logger to find out when the reader is ready.
			if strings.HasPrefix(msg, "Joined group") {
				select {
				case <-joined:
					sc.wg.Done()
				default:
					close(joined) // Close once.
				}
			}
		}),
	})
	select {
	case <-joined:
		sc.wg.Done()
	case <-time.After(10 * time.Second):
		return nil, errors.New("did not join group in time")
	}

	sc.wg.Wait()

	return sc, nil
}

func (sc *SubscriptionClient) Receive(ctx context.Context) (kafka.Message, error) {
	return sc.reader.FetchMessage(ctx)
}

func (sc *SubscriptionClient) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return sc.reader.CommitMessages(ctx, msg)
}
