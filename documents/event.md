# Events
- An Event represents an action that has completed.
- Always use verbs in the past tense: for example SeatReserved.
- Events are raised from within an aggregate, often as a consequence of a command. For example: the command UpdateAddress may result in an event AddressUpdated.
- We talk about events being Published; there may be multiple subscribers to an event.
- Events may be handled by subscribers within the same bounded context, or by another bounded context. If they are handled by another bounded context then there must be some mechanism (such as an event bus) to transport the event. The event will also need to pass through an [Anti corruption layer]().
- Event subscribers may be [Aggregates](#Aggregates) or [Sagas](#Sagas).
- In the [Event sourcing](#Eventsourcing) approach, events are captured and stored. By replaying the stored events, the state of an aggregate can be recreated.
When implementing the CQRS pattern, events may be used as the mechanism for synchronizing the read-side with the write-side. The read-side is a subscriber to the events raised on the write-side.

```golang
    // cqrs event
    type OrderPlaced struct {
        messaging.Event // (*)
        ConferenceId string
	    AccessCode   string
	    Seats        []int
    }

    // eventsourcing event
    type OrderPlaced struct {
	    VersionedEvent // (*)
	    ConferenceId string
	    AccessCode   string
	    Seats        []int
    }
```