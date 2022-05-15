package tuples

import (
	"reflect"
	"strings"
	"testing"
)

type scanTest struct {
	desc string
	in   string
	out  [][][]string // [[[key value], .... pairs of key values is a tuple], ....]
	err  error
}

var scanTests = []scanTest{{
	desc: "Single tuple",
	in:   "fname=John,lname=Doe,dob=2000-01-01",
	out:  [][][]string{{{"fname", "John"}, {"lname", "Doe"}, {"dob", "2000-01-01"}}},
}, {
	desc: "Multiple tuples",
	in:   "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	out: [][][]string{
		{{"fname", "John"}, {"lname", "Doe"}, {"dob", "2000-01-01"}},
		{{"fname", "Bob"}, {"lname", "Smith"}, {"dob", "2010-10-10"}},
	},
}}

func TestNext(t *testing.T) {
	for _, tC := range scanTests {
		t.Run(tC.desc, func(t *testing.T) {
			s := newScanner(strings.NewReader(tC.in))
			var out [][][]string
			for {
				tuple, done := s.next()
				out = append(out, tuple)
				if done {
					break
				}
			}

			if tC.err != nil {

			} else {
				for i, tout := range tC.out {
					if !reflect.DeepEqual(out[i], tout) {
						t.Errorf("#%d: scan tuple() output:\ngot  %v\nwant %v", i, out[i], tout)
					}
				}
				if len(out) != len(tC.out) {
					t.Errorf("scan tuple() output length mismatch:\ngot  %d\nwant %d", len(out), len(tC.out))
				}
			}
		})
	}
}
