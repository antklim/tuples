package tuples_test

import (
	"testing"
	"time"

	"github.com/antklim/tuples"
)

type T1 struct {
	Foo string `tuples:"foo"`
	Bar int    `tuples:"baaar"`
}

type T3 struct {
	Name    string `tuples:"fname"`
	Age     int    `tuples:"age"`
	Dob     time.Time
	HasKids bool `tuples:"with_kids"`
}

type T4 struct {
	F1 string `tuples:"fld4"`
	F2 int
	F3 int64 `tuples:"fld5"`
}

type marshalTest struct {
	in  any
	out string
	err error
}

var marshalTests = []marshalTest{
	// {in: T1{Foo: "hey", Bar: 25}, out: "foo=hey,baaar=25"},
	// {in: &T1{Foo: "hey", Bar: 25}, out: "foo=hey,baaar=25"},
	// {
	// 	in: []T3{
	// 		{Name: "Bob", Age: 33, Dob: time.Now(), HasKids: false},
	// 		{Name: "Junior", Age: 43, Dob: time.Now(), HasKids: true},
	// 	},
	// 	out: "fname=Bob,age=33,with_kids=false fname=Junior,age=43,with_kids=true",
	// },
	// {
	// 	in: &[]T3{
	// 		{Name: "Bob", Age: 33, Dob: time.Now(), HasKids: false},
	// 		{Name: "Junior", Age: 43, Dob: time.Now(), HasKids: true},
	// 	},
	// 	out: "fname=Bob,age=33,with_kids=false fname=Junior,age=43,with_kids=true",
	// },
	{
		in: map[string]any{
			"fname": "Smith",
			"age":   23,
		},
		out: "fname=Smith,age=23",
	},
	// {
	// 	in: []any{
	// 		map[string]any{"fld1": int8(1), "fld2": uint64(9)},
	// 		T4{F1: "hehe", F2: 2, F3: int64(44)},
	// 	},
	// 	out: "fld1=1,fld2=9 fld4=hehe,fld5=44",
	// },
}

func TestMarshal(t *testing.T) {
	for tI, tC := range marshalTests {
		got, err := tuples.Marshal(tC.in)
		if !eqErrors(tC.err, err) {
			t.Errorf("#%d: unexpected Marshal() error: \ngot  %v\nwant %v", tI, err, tC.err)
			continue
		}
		if out := string(got); out != tC.out {
			t.Errorf("#%d: Marshal() output:\ngot  %v (len %d)\nwant %v (len %d)",
				tI, out, len(out), tC.out, len(tC.out))
		}
	}
}
