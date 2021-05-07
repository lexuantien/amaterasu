package messagebroker

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestClient(t *testing.T) {
	w := &kafka.Writer{}
	c := &TopicClient{}
	c.SetWriter(w)
}
