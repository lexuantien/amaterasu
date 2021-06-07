package serialization

import (
	"encoding/json"
	"reflect"
)

type JsonSerializer struct {
}

// create new JsonSerialize
func New_JsonSerializer() JsonSerializer {
	return JsonSerializer{}
}

// Serializes an object graph
// obj object to serialize
// []byte easy put it to kafka
// error	nil if success, else
func (js JsonSerializer) Serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// Deserializes an object graph to specific type
// message input
// entry type of output, 01/01/2022 golang support generic, so we must use this
// interface{} output value
// error if fail, else nil
func (js JsonSerializer) Deserialize(message []byte, entry reflect.Type) (interface{}, error) {

	// create default value
	exportVal := reflect.New(entry).Interface()

	js.Deserialize2(message, &exportVal)

	return exportVal, nil
}

func (js JsonSerializer) Deserialize2(message []byte, output interface{}) {
	json.Unmarshal(message, &output)
}
