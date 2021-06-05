package serialization

import "reflect"

// Interface for serializers that can read/write an object graph to a stream.
type ISerializer interface {

	// Serializes an object graph
	// obj object to serialize
	// []byte easy put it to kafka
	// error	nil if success, else
	Serialize(obj interface{}) ([]byte, error)

	// Deserializes an object graph to specific type
	// message input
	// entry type of output, 01/01/2022 golang support generic, so we must use this
	// interface{} output value
	// error if fail, else nil
	Deserialize(message []byte, entry reflect.Type) (interface{}, error)

	Deserialize2(message []byte, output interface{})
}
