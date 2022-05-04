package tuples_test

import (
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

type unmarshalTest struct {
	in  string
	ptr any
	out any
	err error
}

var unmarshalTests = []unmarshalTest{
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
	// TODO: add float test
	// TODO: add array test
	// TODO: add unmarshal to map test
	// TODO: add unmarshal to interface test
	// TODO: add errors tests
	// 	- unmarshal to struct
	// 	- unmarshal to invalid data types - decode string to int for example
	//	- unmarshal to unsupported data types - decode to slices
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
		err := tuples.Unmarshal(in, got.Interface())
		if err != nil {
			t.Fatalf("unexpected Unmarshal() error: %v", err)
		}
		if !reflect.DeepEqual(got.Elem().Interface(), tC.out) {
			t.Fatalf("Unmarshal() output:\ngot %v\nwant %v", got, tC.out)
		}
	}
}
