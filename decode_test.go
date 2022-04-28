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

func TestUnmarshal(t *testing.T) {
	in := []byte("fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10")
	type out struct {
		Name        string `tuples:"fname"`
		Surname     string `tuples:"lname"`
		DateOfBirth string `tuples:"dob"`
	}
	want := []out{{
		Name:        "John",
		Surname:     "Doe",
		DateOfBirth: "2000-01-01",
	}, {
		Name:        "Bob",
		Surname:     "Smith",
		DateOfBirth: "2010-10-10",
	}}

	var got []out
	err := tuples.Unmarshal(in, &got)
	if err != nil {
		t.Fatalf("unexpected Unmarshal() error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unmarshal() output:\ngot %v\nwant %v", got, want)
	}
}

// func TestUnmarshaSimple(t *testing.T) {
// 	t.Skip()
// 	in := []byte("fname=John,lname=Doe,dob=2000-01-01")
// 	type out struct {
// 		Name        string `tuples:"fname"`
// 		Surname     string `tuples:"lname"`
// 		DateOfBirth string `tuples:"dob"`
// 	}
// 	// want := out{
// 	// 	Name:        "John",
// 	// 	Surname:     "Doe",
// 	// 	DateOfBirth: "2000-01-01",
// 	// }

// 	var got out
// 	err := tuples.Unmarshal(in, &got)
// 	if err != nil {
// 		t.Fatalf("unexpected Unmarshal() error: %v", err)
// 	}
// 	fmt.Printf("%+v\n", got)
// }
