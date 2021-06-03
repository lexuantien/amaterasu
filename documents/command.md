# Commands

- Commands are imperatives (e.g. BookSeats, PlaceOrder); they are requests for the system to perform a task or action. For example, "book two places on conference X" or "allocate speaker Y to room Z."
- The goal of a Command is to mutate the state of an application in some way to a new, different state
- Typically, commands result in an [Events](#Events) being raised.
- Commands are usually processed just once, by a single recipient.
- Commands are one-way and can be delivered asynchronously.
- Commands can be created by a client (for example the UI), and dispatched by a Command handler to the correct [Aggregate](#Aggregate) in the domain.
- Commands can be created by a [Sagas](#Sagas), and dispatched by a [Command handler](#Commandhandler) to the correct Aggregate in the domain.

- Deepdive about command in Amaterasu project:

    - Each command has the Id, 2 fuction getter and setter

```golang
type (
	ICommand interface {
		GetId() string
		SetId(string)
	}

	Command struct {
		Id string
	}
)

func (cmd *Command) GetId() string { /*...*/ }

func (cmd *Command) SetId(id string) { /*...*/ }
```
- Example about using command
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


# Commands and Data Transfer Objects
- A typical approach to enabling a user to edit data is to use data transfer objects (DTO): the UI retrieves the data to be edited from the application as a DTO, a user edits the DTO in the UI, the UI sends the modified DTO back to the application, and then the application applies those changes to the data in the database.

- This approach is data-centric and tends to use standard CRUD operations throughout. In the UI, the user performs operations that are essentially CRUD operations on the data in the DTO.

- This is a simple, well understood aproach that works well for many applications. However, for some applications it is more useful if the UI sends commands instead of DTOs back to the application to make changes to the data. Commands are behavior-centric instead of data-centric, directly represent operations in the domain, maybe more intuitive to users, and can capture the user's intent more effectively than DTOs.

- In a typical CQRS implementation, the read-model returns data to the UI as a DTO. The UI then sends a command (not a DTO) to the write-model.