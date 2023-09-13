package tuples

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Unmarshal parses the tuples-encoded data and stores the result in the value
// pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func Unmarshal(data []byte, v any) error {
	var d decoder

	if err := d.init(data); err != nil {
		return err
	}

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

type decoder struct {
	data []byte
	s    *scanner
}

func (d *decoder) init(data []byte) error {
	d.data = data

	err := d.initScanner(bytes.NewReader(data))

	return err
}

func (d *decoder) initScanner(r io.Reader) error {
	s, err := newScanner(r)
	if err != nil {
		return err
	}

	d.s = s

	return nil
}

func (d *decoder) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return d.value(rv)
}

func (d *decoder) value(v reflect.Value) error {
	v = indirect(v)

	switch d.s.state {
	case scanReady:
		// The beginning of scanning, v should be a slice, array or interface.
		if v.IsValid() {
			if err := d.array(v); err != nil {
				return err
			}
		}
	case scanTuple:
		// In the middle of scanning, v should be a struct.
		if v.IsValid() {
			if err := d.object(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *decoder) array(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		break
	case reflect.Interface:
		return d.arrayInterface(v)
	default:
		return &UnmarshalError{Value: "array", Type: v.Type()}
	}

	// TODO(feat): add decoder options that allow addition of the elements to non-empty
	//			 slice or array. Currently elements are overwritten from the start.
	// i := v.Len()
	i := 0
	for d.s.next() {
		if v.Kind() == reflect.Slice {
			// Grow slice if necessary.
			if i >= v.Cap() {
				newcap := v.Cap() + v.Cap()/2 //nolint: gomnd
				if newcap < 4 {               //nolint: gomnd
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
			// Cut off the rest of prefilled slice.
			v.SetLen(i)
		}
	}

	if i == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}

	return nil
}

func (d *decoder) object(v reflect.Value) error {
	flds, err := d.s.tuple()
	if err != nil {
		return err
	}

	for _, fld := range flds {
		tag, val := fld[idxKey], fld[idxVal]
		sf := cachedTypeFields(v.Type())

		if idx, ok := sf.fieldsByTag[tag]; ok {
			if err := set(v.Field(idx), val); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *decoder) arrayInterface(v reflect.Value) error {
	var a = make([]map[string]any, 0)
	var er error

	for d.s.next() {
		if oi, err := d.objectInterface(); err != nil {
			er = err
			break
		} else if oi != nil {
			a = append(a, oi)
		}
	}

	v.Set(reflect.ValueOf(a))

	return er
}

func (d *decoder) objectInterface() (map[string]any, error) {
	m := make(map[string]any)

	flds, err := d.s.tuple()
	if err != nil {
		return nil, err
	}

	for _, fld := range flds {
		m[fld[idxKey]] = fld[idxVal]
	}

	return m, nil
}

// indirect walks down v until it gets to a non-pointer.
// inspired by
//
//	https://github.com/aws/aws-sdk-go/blob/7a3b8d6ddc7199249e6280d6c1839e08213cc48c/service/dynamodb/dynamodbattribute/decode.go#L634
//	https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/encoding/json/decode.go;l=420;drc=refs%2Ftags%2Fgo1.18.1
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	for {
		if v.Kind() != reflect.Pointer {
			break
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
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}

		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}

		v.SetUint(n)
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return &UnmarshalError{Err: err, Value: value, Type: v.Type()}
		}

		v.SetFloat(f)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
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
