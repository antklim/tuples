package tuples

import (
	"reflect"
	"testing"
)

func TestCachedTypeFields(t *testing.T) {
	var out struct {
		Name        string `tuples:"fname"`
		Surname     string `tuples:"lname"`
		DateOfBirth string
		Age         int `tuples:"age"`
	}

	expected := typFields{
		fields: []field{
			{name: "Name", tag: "fname"},
			{name: "Surname", tag: "lname"},
			{name: "Age", tag: "age"},
		},
		fieldsByTag: map[string]int{
			"fname": 0,
			"lname": 1,
			"age":   2,
		},
	}

	got := cachedTypeFields(reflect.TypeOf(out))
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("typeFields() output:\ngot  %v\nwant %v", got, expected)
	}
}
