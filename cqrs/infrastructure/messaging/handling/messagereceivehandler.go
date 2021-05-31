package handling

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OnProcessMessage func(msg messaging.Envelope) bool
type OnCompleteMessage func(context.Context, *kafka.Message) error

const LOOP_DO_AGAIN = 3

func OnMessageReceivedHandler(ctx context.Context, message *kafka.Message, onProcessMessage OnProcessMessage, onCompleteMessage OnCompleteMessage) error {
	msg := messaging.Envelope{} // wrapper class contain transfer data

	// decode envelop class inside message
	decoder := json.NewDecoder(bytes.NewReader(message.Value))
	decoder.UseNumber()
	errUnmarshal := decoder.Decode(&msg) // parse

	if errUnmarshal != nil {
		return errUnmarshal
	}

	// have a lot of onProcessMessage func, ie: onProcessMessage for command,
	// onProcessMessage for event, onProcessMessage for event sourcing
	processSuccess := onProcessMessage(msg)

	// commit msg
	if processSuccess {
		loop := 0
		for {
			// commit message
			if errCommitMsg := onCompleteMessage(ctx, message); errCommitMsg != nil {
				log.Fatalln(errCommitMsg)
				if loop == LOOP_DO_AGAIN { // re-commit
					// TODO use circuit breaker pattern to turn off service
					// TODO add to dead letter to commit again
					return errors.New("commit message fail")
				}
				loop++
			} else { // success
				// TODO add to `undispatch message` table if type event
				return nil
			}
		}
	}

	// TODO use circuit breaker pattern to turn off service
	// TODO add to `dead letter` table to process again
	return errors.New("process message fail")
}
