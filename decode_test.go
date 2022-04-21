package tuples_test

import (
	"testing"

	"github.com/antklim/tuples"
)

func TestUnmarshal(t *testing.T) {
	t.Skip("not implemented")
	in := []byte("fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10")
	type out struct {
		Name        string `tuples:"fname"`
		Surname     string `tuples:"lname"`
		DateOfBirth string `tuples:"dob"`
	}
	// want := []out{{
	// 	Name:        "John",
	// 	Surname:     "Doe",
	// 	DateOfBirth: "2000-01-01",
	// }, {
	// 	Name:        "Bob",
	// 	Surname:     "Smith",
	// 	DateOfBirth: "2010-10-10",
	// }}

	var got []out
	err := tuples.Unmarshal(in, &got)
	if err != nil {
		t.Fatalf("unexpected Unmarshal() error: %v", err)
	}
}
