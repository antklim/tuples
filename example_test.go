package tuples_test

import (
	"fmt"
	"strings"

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
	}
	fmt.Printf("%+v\n", p)

	in = "name=John,age=23,adult=true name=Bob,adult=true adult=true,age=30"
	var pp []person
	if err := tuples.Unmarshal([]byte(in), &pp); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", pp)

	in = "name=John,lname=Doe,age=17 height=170,weight=50"
	var ppp any
	if err := tuples.Unmarshal([]byte(in), &ppp); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ppp)

	// Output:
	// [{Name:John Age:23 IsAdult:true}]
	// [{Name:John Age:23 IsAdult:true} {Name:Bob Age:0 IsAdult:true} {Name: Age:30 IsAdult:true}]
	// [map[age:17 lname:Doe name:John] map[height:170 weight:50]]
}

func ExampleMarshal() {
	type person struct {
		Name string `tuples:"full_name"`
		Age  int    `tuples:"years-old"`
	}

	p := person{Name: "Bob", Age: 33}

	b, err := tuples.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", string(b))

	pp := []person{{Name: "Bob", Age: 33}, {Name: "Paul", Age: 99}}

	b, err = tuples.Marshal(pp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", string(b))

	// Output:
	// full_name=Bob,years-old=33
	// full_name=Bob,years-old=33 full_name=Paul,years-old=99
}

func ExampleReadString() {
	in := "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10"

	v, err := tuples.ReadString(in)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", v)

	// Output:
	// [[John Doe 2000-01-01] [Bob Smith 2010-10-10]]
}

func ExampleReader_ReadAll() {
	in := "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10"

	r, err := tuples.NewReader(strings.NewReader(in))
	if err != nil {
		fmt.Println(err)
	}

	v, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", v)

	// Output:
	// [[John Doe 2000-01-01] [Bob Smith 2010-10-10]]
}
