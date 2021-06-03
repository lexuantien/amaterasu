# Commandhandler 
- A Command handler is responsible for dispatching the command to the correct domain object Aggregate root.
- A command handler checks that the command is correctly formatted ('sanity check'). This check is separate from the Anti-corruption Layer, and verifies that the state of the system is compatible with command execution.
- A command handler locates the aggregate instance that will action (provide an execution context for) the command. This may require instantiating a new instance or rehydrating an existing instance.
- A command handler invokes a business method of the aggregate passing any parameters from the command.
- The aggregate executes the business logic and creates any events inside the business method.
- After the business method completes, the command handler persists the aggregate's changes to the repository.
- In Event sourcing, the command handler appends any new events to the event store.
- The events are published to a topic when the aggregate state is persisted (saving the state and publishing an event may be wrapped inside a transaction).


# Deep Dive

- Using command processor to process command. Command handler must register to processor to get command

```golang
	consumer := messaging.New_Consumer(config)
	serializer := serialization.New_JsonSerializer()

	handler1 := New_CommandHander1()
	processor := handling.New_CommandProcessor(/* some config */)

    // register command handler to processor
	processor.Register(handler1)
```

- In command processor using command dispatcher to store command handler

```golang
func (cd *CommandDispatcher) Register(commandHandler ICommandHandler) error {

	handlerType := reflect.TypeOf(commandHandler)
	numHandlerMethods := handlerType.NumMethod() // count all method

    // ...

	for i := 0; i < numHandlerMethods; i++ {
		// get a method in commandHandler
		handlerMethod := handlerType.Method(i)

        // ...

		// func (fh fooCommandHandler) handleFoo1(f foo1) error {}
		// func (fh fooCommandHandler) handleFoo2(f foo2) error {}
		handlerFunc := func(msg messaging.Message) error {

			response := handlerMethod.Func.Call([]reflect.Value{
				reflect.ValueOf(commandHandler), // fh param
				reflect.ValueOf(msg).Elem(),     // f param
			})

			return nil
		}

		// 1 because param in golang start at 1, 0 is struct type
		commandType := handlerMethod.Type.In(1) // get command type

		// 1 command is handled only by 1 commandhandler, failed if 2 or 3 ...
		if _, ok := cd.commandHandlers[commandType]; ok {
			return errors.New("the command handled by the received handler already has a registered handler")
		}

		// command type as string
		commandTypeName := utils.GetObjType2(commandType)
		// store
		cd.commandHandlers[commandType] = commandHandler
		cd.commandTypes[commandTypeName] = commandType
		cd.commandHandlerFuncs[commandTypeName] = handlerFunc

	}

	return nil
}

``` 

- Example using command:
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
