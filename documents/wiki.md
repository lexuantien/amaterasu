

# Commands
- Commands are imperatives (e.g. BookSeats, PlaceOrder); they are requests for the system to perform a task or action. For example, "book two places on conference X" or "allocate speaker Y to room Z."
- The goal of a Command is to mutate the state of an application in some way to a new, different state
- Typically, commands result in an [Events](#Events) being raised.
- Commands are usually processed just once, by a single recipient.
- Commands are one-way and can be delivered asynchronously.
- Commands can be created by a client (for example the UI), and dispatched by a Command handler to the correct [Aggregate](#Aggregate) in the domain.
- Commands can be created by a [Sagas](#Sagas), and dispatched by a [Command handler](#Commandhandler) to the correct Aggregate in the domain.
- Can an aggregate create and send a command? No
- Is a command a DTO? No. Yes if you don’t need backward compatibility
```golang
    // cqrs command
    type ConfirmOrder struct {        
        messaging.Command // (*)
        OrderId string
    }

    type AddSeats struct {
        messaging.Command  // (*)
        SeatAssignmentsId string
        Position          int
        Attendee          string
    }
```
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
# Commandhandler

- A Command handler is responsible for dispatching the command to the correct domain object Aggregate root.
- A command handler checks that the command is correctly formatted ('sanity check'). This check is separate from the Anti-corruption Layer, and verifies that the state of the system is compatible with command execution.
- A command handler locates the aggregate instance that will action (provide an execution context for) the command. This may require instantiating a new instance or rehydrating an existing instance.
- A command handler invokes a business method of the aggregate passing any parameters from the command.
- The aggregate executes the business logic and creates any events inside the business method.
- After the business method completes, the command handler persists the aggregate's changes to the repository.
- In Event sourcing, the command handler appends any new events to the event store.
- The events are published to a topic when the aggregate state is persisted (saving the state and publishing an event may be wrapped inside a transaction).

```golang
    // create handler 
    type OrderCommandHander struct{
        dbContext *eventsourcing.EventStore
    }
    
    func (o CommandHander1) HandleConfirmOrder(command ConfirmOrder) error {

        // find order aggreate in event_store table
        order := dbContext.Find(command.OrderId).(*Order)

        if order != nil {
            order.Update(&OrderConfirmed{/* some data here */ })
        }

        // save an event
	    dbContext.Save(order)
	    return nil
    }

    func (o CommandHander1) HandleAddSeats(command AddSeats) error {
	    fmt.Println(command)
	    return nil
    }

```

# Eventhandler

Events are published to multiple recipients, typically aggregate instances or process managers. The Event handler performs the following tasks:
- It receives an Event instance from the messaging infrastructure.
- It locates the aggregate or process manager instance that is the target of the event. This may involve creating a new aggregate instance or locating an existing instance.
- It invokes the appropriate method on the aggregate or process manager instance, passing in any parameters from the event.
- It persists the new state of the aggregate or process manager to storage.

# Eventsourcing

# Idempotent
An operation that when repeated gives the same result. This concept is used in mathematics as well; but becomes an important concept in messaging when we deal with things like guaranteed delivery or at-least-once delivery.

# Locking
- These are stategies to handle contention for resources.
- There are two general models to deal with these situations that you can observe in the wild:

    - Pessimistic. You take a time-bombed lock on an object and reserve it for the client. Ticketmaster operates their web site that way. When they show you a set of seats tickets they are yours until the clock ticks down.
    - Optimistic. There's no lock. Expedia/Orbitz work that way. When they show you a seat, you don't get it until you book it and there is no promise made.

# Snapshot
A term related to event sourcing. It refers to an optimization when you replay an event stream, whereby you only need to replay the stream from the last snapshot rather than replaying all events.

