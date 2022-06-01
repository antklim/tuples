package tuples_test

import (
	"testing"
	"time"
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

var in1 = T1{Foo: "hey", Bar: 25}
var out1 = "foo=hey,baaar=25"

var in2 = []T3{
	{Name: "Bob", Age: 33, Dob: time.Now(), HasKids: false},
	{Name: "Junior", Age: 43, Dob: time.Now(), HasKids: true},
}
var out2 = "fname=Bob,age=33,with_kids=false fname=Junior,age=43,with_kids=true"

var in3 = map[string]any{
	"fname": "Smith",
	"age":   23,
}
var out3 = "fname=Smith,age=23"

var in4 = []map[string]any{
	{"foo": "bar", "temp": uint(32), "Boool": true},
	{"price": 15.5, "fee": float64(3.4), "tempo": "jazz"},
}
var out4 = "foo=bar,temp=32,Boool=true price=15.5,fee=3.4,tempo=jazz"

var in5 = []any{
	map[string]any{"fld1": int8(1), "fld2": uint64(9)},
	T4{F1: "hehe", F2: 2, F3: int64(44)},
}
var out5 = "fld1=1,fld2=9 fld4=hehe,fld5=44"

func TestMarshal(t *testing.T) {
	// for _, tC := range testCases {
	// 	t.Run(tC.desc, func(t *testing.T) {

	// 	})
	// }
}
