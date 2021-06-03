# Eventhandler
Events are published to multiple recipients, typically aggregate instances or process managers. The Event handler performs the following tasks:
- It receives an Event instance from the messaging infrastructure.
- It locates the aggregate or process manager instance that is the target of the event. This may involve creating a new aggregate instance or locating an existing instance.
- It invokes the appropriate method on the aggregate or process manager instance, passing in any parameters from the event.
- It persists the new state of the aggregate or process manager to storage.
