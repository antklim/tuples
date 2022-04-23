package tuples

import (
	"errors"
	"reflect"
)

// Unmarshal parses the tuples-encoded data and stores the result in the value
// pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func Unmarshal(data []byte, v any) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}

// InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer).
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "tuples: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "tuples: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "tuples: Unmarshal(nil " + e.Type.String() + ")"
}

type decodeState struct {
	data []byte
}

func (d *decodeState) init(data []byte) {
	d.data = data
}

func (d *decodeState) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return errors.New("not implemented")
}
