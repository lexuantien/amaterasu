package serialization

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type JsonSerializer struct{}

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
func (js JsonSerializer) Deserialize(message interface{}, entry reflect.Type) (interface{}, error) {

	// create default value
	exportVal := reflect.New(entry).Interface()

	// config map data
	config := &mapstructure.DecoderConfig{
		// DecodeHook: utils.MapTimeFromJSON,
		// TagName:          "json",
		Result:   exportVal,
		Metadata: nil,
		// WeaklyTypedInput: true,
	}

	decoder, errDecoder := mapstructure.NewDecoder(config)
	if errDecoder != nil {
		return nil, errors.New("config mapstructure fail")
	}

	errDecode := decoder.Decode(message)
	if errDecode != nil {
		return nil, errors.New("decode message fail")
	}

	// json.Unmarshal()

	return exportVal, nil
}
