package handling

// Provides basic common processing code for components that handle
// incoming messages from a receiver.
type IMessageProcessor interface {
	Start() // Starts the listener.
	Stop()  // Stops the listener.
}
