package tuples_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/antklim/tuples"
)

type readTest struct {
	Name   string
	Input  string
	Output [][]string
}

var readTests = []readTest{{
	Name:   "Simple",
	Input:  "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}}

func TestRead(t *testing.T) {
	for _, tC := range readTests {
		t.Run(tC.Name, func(t *testing.T) {
			r := tuples.NewReader(strings.NewReader(tC.Input))
			out, err := r.ReadAll()
			if err != nil {
				t.Fatalf("unexpected ReadAll() error: %v", err)
			}
			if !reflect.DeepEqual(out, tC.Output) {
				t.Fatalf("ReadAll() output:\ngot %v\nwant %v", out, tC.Output)
			}

			r = tuples.NewReader(strings.NewReader(tC.Input))
			for recNum := 0; ; recNum++ {
				rec, err := r.Read()
				var wantErr error
				if recNum >= len(tC.Output) {
					wantErr = io.EOF
				}
				if err != nil && err.Error() != wantErr.Error() {
					t.Fatalf("Read() error at record %d:\ngot %v\nwant %v", recNum, err, wantErr)
				}
				if err != nil {
					break
				}
				if got, want := rec, tC.Output[recNum]; !reflect.DeepEqual(got, want) {
					t.Errorf("Read vs ReadAll mismatch:\ngot %v\nwant %v", got, want)
				}
			}
		})
	}
}
