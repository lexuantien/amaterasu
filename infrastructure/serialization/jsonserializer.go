package serialization

import "encoding/json"

type JsonSerializer struct{}

// :Create
func New_JsonSerializer() JsonSerializer {
	return JsonSerializer{}
}

// Serialize data
func (js JsonSerializer) Serialize(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "")
}

// deserialize data
func (js JsonSerializer) Deserialize(message []byte) (interface{}, error) {

	return struct{}{}, nil
}
