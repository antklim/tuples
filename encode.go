package tuples

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

const (
	tuplesDelimiter = ' '
	fieldsDelimiter = ','
	keyValDelimiter = '='
)

type MarshalError struct {
	err error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("tuples: marshal failed: %s", e.err)
}

func (e *MarshalError) Unwrap() error {
	return e.err
}

// Marshal returns tuples encoding of v.
//
// Marshal supports:
//	- struct with tags
//	- pointer to a struct with tags
//	- slice of structs with tags
//	- pointer to slice of structs with tags
//	- map
//	- pointer to a map
//	- slice of maps
//	- pointer to a slice of maps
//
// In case of map, map key value used as tuple key. All map entries marshaled.
//
// Only basic types supported as values, i.e string, int, float, boolean.
// MarshalError returned in case, when unsupported type found.
func Marshal(v any) ([]byte, error) {
	var e encoder
	if err := e.encode(reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return e.b.Bytes(), nil
}

type encoder struct {
	b bytes.Buffer
}

func (e *encoder) encode(v reflect.Value) error {
	v = unwrapElement(v)

	switch v.Kind() {
	case reflect.Struct:
		return e.structObj(v)
	case reflect.Map:
		return e.mapObj(v)
	case reflect.Slice, reflect.Array:
		return e.array(v)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return nil
	default:
		return e.value(v)
	}
}

func (e *encoder) structObj(v reflect.Value) error {
	return errors.New("not implemented")
}

func (e *encoder) mapObj(v reflect.Value) error {
	for i, key := range v.MapKeys() {
		mapKey := fmt.Sprint(key.Interface())
		if mapKey == "" {
			return &MarshalError{errors.New("map key cannot be empty")}
		}

		if err := e.writeKey(mapKey, i); err != nil {
			return &MarshalError{err}
		}

		keyVal := v.MapIndex(key)
		if err := e.encode(keyVal); err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) array(v reflect.Value) error {
	return errors.New("not implemented")
}

func (e *encoder) value(v reflect.Value) error {
	elem := fmt.Sprint(v.Interface())
	if _, err := e.b.WriteString(elem); err != nil {
		return &MarshalError{err}
	}
	return nil
}

func (e *encoder) writeKey(key string, keyIdx int) error {
	if keyIdx > 0 {
		if _, err := e.b.WriteRune(fieldsDelimiter); err != nil {
			return err
		}
	}

	if _, err := e.b.WriteString(key); err != nil {
		return err
	}

	if _, err := e.b.WriteRune(keyValDelimiter); err != nil {
		return err
	}

	return nil
}

func unwrapElement(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}