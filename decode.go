package tuples

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
		return fmt.Sprintf("tuples: Unmarshal(non-pointer %s)", e.Type.String())
	}
	return fmt.Sprintf("tuples: Unmarshal(nil %s)", e.Type.String())
}

// UnmarshalUnsupportedTypeError describes an unsupported field type of a
// value of a specific Go type.
type UnmarshalUnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnmarshalUnsupportedTypeError) Error() string {
	return fmt.Sprintf("tuples: unsupported Go value type %s",
		e.Type.String())
}

// UnmarshalError describes an error that occurred while unmarshaling a tuple
// fields values into a Go type fields.
type UnmarshalError struct {
	Err   error
	Value string
	Type  reflect.Type
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("tuples: cannot unmarshal %q into Go value of type %s",
		e.Value, e.Type.String())
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

type decodeState struct {
	data   []byte
	opcode int
	s      *bufio.Scanner // TODO: create a separate scanner and reuse it in reader
}

func (d *decodeState) init(data []byte) {
	d.data = data
	d.s = bufio.NewScanner(bytes.NewReader(data))
	d.s.Split(bufio.ScanWords)
}

func (d *decodeState) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return d.value(rv)
}

func (d *decodeState) value(v reflect.Value) error {
	v = indirect(v)
	// TODO: add opcode constants
	switch d.opcode {
	case 0:
		// if it's the beginning of decode then v should be slice, array or interface
		d.opcode = 1
		if v.IsValid() {
			if err := d.array(v); err != nil {
				return err
			}
		}
	case 1:
		// if it's an iteration of the decoding then v should be a struct
		if v.IsValid() {
			if err := d.object(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *decodeState) array(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		break
	case reflect.Interface:
		return d.arrayInterface(v)
	default:
		return &UnmarshalError{Value: "array", Type: v.Type()}
	}

	// TODO: add decoder options that allow addition of the elements to non-empty
	//			 slice or array. Currently elements are overwritten from the start.
	// i := v.Len()
	i := 0
	for d.s.Scan() {
		if v.Kind() == reflect.Slice {
			// Grow slice if necessary.
			if i >= v.Cap() {
				newcap := v.Cap() + v.Cap()/2 // nolint: gomnd
				if newcap < 4 {               // nolint: gomnd
					newcap = 4
				}
				newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
				reflect.Copy(newv, v)
				v.Set(newv)
			}
			if i >= v.Len() {
				v.SetLen(i + 1)
			}
		}

		if i < v.Len() {
			// Decode into element.
			if err := d.value(v.Index(i)); err != nil {
				return err
			}
		} else {
			// Run out of fixed array, skip.
		}
		i++
	}

	if i < v.Len() {
		if v.Kind() == reflect.Array {
			// Add zeros to the rest of the array.
			z := reflect.Zero(v.Type().Elem())
			for ; i < v.Len(); i++ {
				v.Index(i).Set(z)
			}
		} else {
			v.SetLen(i) // ???? is it in case when slice is prefilled ???? add test
		}
	}

	if i == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
	return nil
}

func (d *decodeState) object(v reflect.Value) error {
	t := d.s.Text()
	fv := readFields(t)

	for _, fld := range fv {
		fname, fvalue := fld[0], fld[1]
		for i := 0; i < v.Type().NumField(); i++ {
			field := v.Type().Field(i)
			tag := field.Tag.Get("tuples")
			if tag == fname {
				if err := set(v.Field(i), fvalue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d *decodeState) arrayInterface(v reflect.Value) error {
	var a = make([]map[string]any, 0)
	for d.s.Scan() {
		a = append(a, d.objectInterface())
	}
	v.Set(reflect.ValueOf(a))
	return nil
}

func (d *decodeState) objectInterface() map[string]any {
	m := make(map[string]any)
	t := d.s.Text()
	fv := readFields(t)
	for _, fld := range fv {
		m[fld[0]] = fld[1]
	}
	return m
}

// readFields reads a raw tuple string and returns a slice of tuple's fields.
// Every field represented by field name, field value pair. For example:
//
//	s := "name=John,lname=Doe,age=17"
//	fmt.Println(readFields(s)) // [[name John] [lname Doe] [age 17]]
func readFields(s string) [][]string { // TODO: can return [][2]string
	var fieldValues [][]string
	fields := strings.FieldsFunc(s, func(c rune) bool { return c == ',' })
	for _, f := range fields {
		kv := strings.FieldsFunc(f, func(c rune) bool { return c == '=' })
		fieldValues = append(fieldValues, kv)
	}
	return fieldValues
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// If it encounters an Unmarshaler, indirect stops and returns that.
// If decodingNull is true, indirect stops at the first settable pointer so it
// can be set to nil.
// original https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/encoding/json/decode.go;l=420;drc=refs%2Ftags%2Fgo1.18.1
// and https://github.com/aws/aws-sdk-go/blob/7a3b8d6ddc7199249e6280d6c1839e08213cc48c/service/dynamodb/dynamodbattribute/decode.go#L634
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Pointer && !e.IsNil() {
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Pointer {
			break
		}
		if v.Elem().Kind() != reflect.Pointer && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func set(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64) // nolint:gomnd
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(value, 10, 64) // nolint:gomnd
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}
		v.SetUint(n)
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32) // nolint:gomnd
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}
		v.SetFloat(f)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64) // nolint:gomnd
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}
		v.SetBool(b)
	default:
		return &UnmarshalUnsupportedTypeError{v.Type()}
	}
	return nil
}
