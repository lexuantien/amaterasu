package serialization

import "encoding/json"

type JsonSerializer struct {
}

func (js JsonSerializer) Serialize(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "")
}
