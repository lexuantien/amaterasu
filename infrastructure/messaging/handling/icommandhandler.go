package handling

import "leech-service/infrastructure/messaging"

type ICommangHandler interface {
	Handle(messaging.Command) error
}
