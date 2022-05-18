package tuples_test

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/antklim/tuples"
)

var invalidUnmarshalTests = []struct {
	v   any
	err string
}{
	{nil, "tuples: Unmarshal(nil)"},
	{struct{}{}, "tuples: Unmarshal(non-pointer struct {})"},
	{(*int)(nil), "tuples: Unmarshal(nil *int)"},
}

func TestInvalidUnmarshal(t *testing.T) {
	in := []byte("fname=John")
	for _, tC := range invalidUnmarshalTests {
		err := tuples.Unmarshal(in, tC.v)
		if err == nil || (err.Error() != tC.err) {
			t.Errorf("Unmarshal() error mismatch:\ngot %v\nwant %s", err, tC.err)
		}
	}
}

type T struct {
	Name    string `tuples:"name"`
	Age     int    `tuples:"age"`
	IsAdult bool   `tuples:"adult"`
}

type TInts struct {
	N   int   `tuples:"n"`
	N8  int8  `tuples:"n8"`
	N16 int16 `tuples:"n16"`
	N32 int32 `tuples:"n32"`
	N64 int64 `tuples:"n64"`

	UN   uint   `tuples:"un"`
	UN8  uint8  `tuples:"un8"`
	UN16 uint16 `tuples:"un16"`
	UN32 uint32 `tuples:"un32"`
	UN64 uint64 `tuples:"un64"`
}

type TFloats struct {
	F32 float32 `tuples:"f32"`
	F64 float64 `tuples:"f64"`
}

type TUnsupportedFldType struct {
	A []string `tuples:"a"`
}

type T2 struct {
	Name string
}

func (t2 T2) String() string { return t2.Name }

type UnmarshalTest struct {
	in         string
	ptr        any
	out        any
	err        error
	withUnwrap bool
}

var unmarshalTests = []UnmarshalTest{
	{in: "", ptr: new([]T), out: []T{}},
	{
		in:  "name=John,lname=Doe,age=17",
		ptr: new([]T),
		out: []T{{Name: "John", Age: 17}},
	},
	{
		in:  "name=John,age=23,adult=true",
		ptr: new([]T),
		out: []T{{Name: "John", Age: 23, IsAdult: true}},
	},
	{
		in:  "name=John,age=23,adult=true name=Bob,adult=true adult=true,age=30",
		ptr: new([]T),
		out: []T{
			{Name: "John", Age: 23, IsAdult: true},
			{Name: "Bob", IsAdult: true},
			{Age: 30, IsAdult: true},
		},
	},

	// add only the first record
	{
		in:  "name=John,age=23,adult=true name=Bob,adult=true",
		ptr: new([1]T),
		out: [1]T{{Name: "John", Age: 23, IsAdult: true}},
	},

	// fill the rest of the array with zeros
	{
		in:  "name=John,age=23,adult=true",
		ptr: new([2]T),
		out: [2]T{{Name: "John", Age: 23, IsAdult: true}, {}},
	},

	// unmarshal numbers
	{
		in:  "n=1,n8=-2,n16=3,n32=-4,n64=5 un=11,un8=12,un16=13,un32=14,un64=15 n=21,un8=22,n16=23,un32=24,n64=25",
		ptr: new([]TInts),
		out: []TInts{
			{N: 1, N8: -2, N16: 3, N32: -4, N64: 5},
			{UN: 11, UN8: 12, UN16: 13, UN32: 14, UN64: 15},
			{N: 21, UN8: 22, N16: 23, UN32: 24, N64: 25},
		},
	},
	{
		in:  "f32=1 f64=2 f32=3.4,f64=5.6 f32=inf,f64=inf f32=-inf,f64=-inf f32=+Inf,f64=+Inf f32=-Inf,f64=-Inf",
		ptr: new([]TFloats),
		out: []TFloats{
			{F32: 1.0},
			{F64: 2.0},
			{F32: 3.4, F64: 5.6},
			{F32: float32(math.Inf(0)), F64: math.Inf(0)},
			{F32: float32(math.Inf(-1)), F64: math.Inf(-1)},
			{F32: float32(math.Inf(0)), F64: math.Inf(0)},
			{F32: float32(math.Inf(-1)), F64: math.Inf(-1)},
		},
	},

	// unmarshal to struct
	{
		in:  "name=John,lname=Doe,age=17",
		ptr: new(T),
		err: &tuples.UnmarshalError{Value: "array", Type: reflect.TypeOf(T{})},
	},

	// unmarshal to interface
	{
		in:  "name=John,lname=Doe,age=17 n=1,n8=-2,n16=3,n32=-4,n64=5",
		ptr: new(any),
		out: []map[string]any{
			{"name": "John", "lname": "Doe", "age": "17"},
			{"n": "1", "n8": "-2", "n16": "3", "n32": "-4", "n64": "5"},
		},
	},

	// invalid field value errors
	{
		in:         "adult=a",
		ptr:        new([]T),
		err:        &tuples.UnmarshalError{Value: "a", Type: reflect.TypeOf(true)},
		withUnwrap: true,
	},
	{
		in:         "n=a",
		ptr:        new([]TInts),
		err:        &tuples.UnmarshalError{Value: "a", Type: reflect.TypeOf(1)},
		withUnwrap: true,
	},
	{
		in:         "un=a",
		ptr:        new([]TInts),
		err:        &tuples.UnmarshalError{Value: "a", Type: reflect.TypeOf(uint(1))},
		withUnwrap: true,
	},
	{
		in:         "f32=a",
		ptr:        new([]TFloats),
		err:        &tuples.UnmarshalError{Value: "a", Type: reflect.TypeOf(float32(1))},
		withUnwrap: true,
	},
	{
		in:         "f64=a",
		ptr:        new([]TFloats),
		err:        &tuples.UnmarshalError{Value: "a", Type: reflect.TypeOf(float64(1))},
		withUnwrap: true,
	},

	// unsupported field type error
	{
		in:  "a=a",
		ptr: new([]TUnsupportedFldType),
		err: &tuples.UnmarshalUnsupportedTypeError{Type: reflect.TypeOf([]string{})},
	},

	// invalid tuple expression
	{
		in:  "a=a,b",
		ptr: new([]T),
		err: errors.New("tuples: tuple #1 invalid field #2"),
	},

	// invalid tuple expression and unmarshal to interface
	{
		in:  "name=John,lname=Doe,age=17,1",
		ptr: new(any),
		err: errors.New("tuples: tuple #1 invalid field #4"),
	},
}

func eqErrors(a, b error) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}
	return a.Error() == b.Error()
}

func TestUnmarshal(t *testing.T) {
	for i, tC := range unmarshalTests {
		in := []byte(tC.in)
		typ := reflect.TypeOf(tC.ptr)
		if typ.Kind() != reflect.Pointer {
			t.Errorf("#%d: unmarshalTest.ptr %T is not a pointer type", i, tC.ptr)
			continue
		}
		typ = typ.Elem()
		got := reflect.New(typ)
		if err := tuples.Unmarshal(in, got.Interface()); !eqErrors(err, tC.err) {
			t.Errorf("#%d: unexpected Unmarshal() error: \ngot  %v\nwant %v", i, err, tC.err)
			continue
		} else if err != nil {
			if tC.withUnwrap && errors.Unwrap(err) == nil {
				t.Errorf("#%d: Unmarshal() error should wrap original error", i)
			}
			continue
		}

		if !reflect.DeepEqual(got.Elem().Interface(), tC.out) {
			t.Errorf("#%d: Unmarshal() output:\ngot  %v\nwant %v", i, got.Elem().Interface(), tC.out)
			continue
		}
	}
}

func TestUnmarshalToPrefilledSlice(t *testing.T) {
	cp, ln := 5, 1
	got := make([]T, cp)
	in := "name=John,lname=Doe,age=17"
	err := tuples.Unmarshal([]byte(in), &got)
	if err != nil {
		t.Fatalf("unexpected Unmarshal() error: \ngot %v", err)
	}
	if len(got) != ln {
		t.Errorf("Unmarshal() output len: \ngot  %d\nwant %d", len(got), ln)
	}
	if cap(got) != cp {
		t.Errorf("Unmarshal() output cap: \ngot  %d\nwant %d", len(got), cp)
	}
}
