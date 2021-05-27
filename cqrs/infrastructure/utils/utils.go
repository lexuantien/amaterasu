package utils

import (
	"reflect"
	"strings"
	"time"
)

func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)

	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	// we need to extract only the name without the package
	// name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}

func GetTypeName2(t reflect.Type) string {
	segments := strings.Split(t.String(), ".")
	return segments[len(segments)-1]
}

// MapTimeFromJSON is a decoder hook that maps time data from JSON values, avoiding the issue
// of things appearing as errors/blank when dealing with native Go time types. This is based on
// the code at https://github.com/mitchellh/mapstructure/issues/41
func MapTimeFromJSON(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}
