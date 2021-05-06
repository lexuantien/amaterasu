package msgqueue

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestClient(t *testing.T) {
	w := &kafka.Writer{}
	c := &Client{}
	c.SetWriter(w)
}
