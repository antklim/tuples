package tuples

import (
	"errors"
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
	desc: "Empty input",
}, {
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
}, {
	desc: "Skip empty field",
	in:   "fname=John,,dob=2000-01-01",
	out:  [][][]string{{{"fname", "John"}, {"dob", "2000-01-01"}}},
}, {
	desc: "Invalid field #1",
	in:   "fname=John,lname=Doe,dob=2000-01-01 name,lname=Smith,dob=2010-10-10",
	err:  errors.New("tuples: tuple #2 invalid field #1"),
}, {
	desc: "Invalid field #2",
	in:   "fname=John,lname=Doe,dob=2000-01-01 name=,lname=Smith,dob=2010-10-10",
	err:  errors.New("tuples: tuple #2 invalid field #1"),
}, {
	desc: "Invalid field #3",
	in:   "fname=John,lname=Doe,dob=2000-01-01 =Bob,lname=Smith,dob=2010-10-10",
	err:  errors.New("tuples: tuple #2 invalid field #1"),
}}

func TestNext(t *testing.T) {
	for _, tC := range scanTests {
		t.Run(tC.desc, func(t *testing.T) {
			s := newScanner(strings.NewReader(tC.in))
			var out [][][]string
			for {
				tuple, done := s.next()
				if tuple != nil {
					out = append(out, tuple)
				}
				if done {
					break
				}
			}

			if tC.err != nil {
				if s.err == nil || (s.err.Error() != tC.err.Error()) {
					t.Fatalf("scan next() error mismatch:\ngot  %v\nwant %v", s.err, tC.err)
				}
			} else {
				for i, tout := range tC.out {
					if !reflect.DeepEqual(out[i], tout) {
						t.Errorf("#%d: scan next() output:\ngot  %v\nwant %v", i, out[i], tout)
					}
				}
				if len(out) != len(tC.out) {
					t.Errorf("scan next() output length mismatch:\ngot  %d\nwant %d", len(out), len(tC.out))
				}
			}
		})
	}
}

func TestNextAfterDone(t *testing.T) {
	s := newScanner(strings.NewReader("fname=John"))
	s.next()
	out, done := s.next()
	if out != nil {
		t.Errorf("scan next() output:\ngot  %v\nwant nil", out)
	}
	if !done {
		t.Errorf("scan next() done:\ngot  %t\nwant true", done)
	}
}
