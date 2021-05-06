package serialization

// Interface for serializers that can read/write an object graph to a stream.
type ISerializer interface {
	// Serializes an object graph
	Serialize(obj interface{}) ([]byte, error)
	// Deserializes an object graph from the specified text reader.
	Deserialize() interface{}
}
