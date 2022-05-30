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

	expected := structFields{
		fields: []field{
			{name: "Name", typ: reflect.TypeOf(""), tag: true},
			{name: "Surname", typ: reflect.TypeOf(""), tag: true},
			{name: "Age", typ: reflect.TypeOf(1), tag: true},
		},
		namedIdx: map[string]int{
			"Name":    0,
			"Surname": 1,
			"Age":     2,
		},
	}

	got := cachedTypeFields(reflect.TypeOf(out))
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("typeFields() output:\ngot  %v\nwant %v", got, expected)
	}
}
