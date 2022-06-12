package tuples_test

import (
	"fmt"

	"github.com/antklim/tuples"
)

func ExampleUnmarshal() {
	type person struct {
		Name    string `tuples:"name"`
		Age     int    `tuples:"age"`
		IsAdult bool   `tuples:"adult"`
	}

	in := "name=John,age=23,adult=true"
	var p []person
	if err := tuples.Unmarshal([]byte(in), &p); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", p)

	in = "name=John,age=23,adult=true name=Bob,adult=true adult=true,age=30"
	var pp []person
	if err := tuples.Unmarshal([]byte(in), &pp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", pp)

	in = "name=John,lname=Doe,age=17 height=170,weight=50"
	var ppp any
	if err := tuples.Unmarshal([]byte(in), &ppp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", ppp)

	// Output:
	// [{Name:John Age:23 IsAdult:true}]
	// [{Name:John Age:23 IsAdult:true} {Name:Bob Age:0 IsAdult:true} {Name: Age:30 IsAdult:true}]
	// [map[age:17 lname:Doe name:John] map[height:170 weight:50]]
}

func ExampleMarshal() {}

func ExampleReadString() {}

func ExampleReader_ReadAll() {}
