package tuples_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/antklim/tuples"
)

type readTest struct {
	desc string
	in   string
	out  [][]string
	err  error

	fDelim  rune
	kvDelim rune
}

var errInvalidDelim = errors.New("tuples: invalid fields or key values delimiter")

var readTests = []readTest{{
	desc: "Simple",
	in:   "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	out:  [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	desc: "EOLTest",
	in:   "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10\n",
	out:  [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	desc: "TrailingWhitespace",
	in:   "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10   ",
	out:  [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	desc: "MultipleWhitespaces",
	in:   "fname=John,lname=Doe,dob=2000-01-01    fname=Bob,lname=Smith,dob=2010-10-10\n",
	out:  [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	desc: "TabDelimited",
	in: "fname=John,lname=Doe,dob=2000-01-01	fname=Bob,lname=Smith,dob=2010-10-10\n",
	out: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	desc:    "CustomDelimiters",
	in:      "fname:John;lname:Doe;dob:2000-01-01 fname:Bob;lname:Smith;dob:2010-10-10",
	out:     [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
	fDelim:  ';',
	kvDelim: ':',
}, {
	desc:    "BadDelimiters1",
	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	err:     errInvalidDelim,
	fDelim:  ':',
	kvDelim: ':',
}, {
	desc:    "BadDelimiters2",
	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	err:     errInvalidDelim,
	fDelim:  utf8.RuneError,
	kvDelim: ':',
}, {
	desc:    "BadDelimiters3",
	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	err:     errInvalidDelim,
	fDelim:  ':',
	kvDelim: utf8.RuneError,
}}

func newReader(rt readTest) *tuples.Reader {
	r := tuples.NewReader(strings.NewReader(rt.in))
	if rt.fDelim != 0 {
		r.FieldsDelimiter = rt.fDelim
	}
	if rt.kvDelim != 0 {
		r.KeyValDelimiter = rt.kvDelim
	}
	return r
}

func TestRead(t *testing.T) {
	for _, tC := range readTests {
		t.Run(tC.desc, func(t *testing.T) {
			r := newReader(tC)
			out, err := r.ReadAll()

			if tC.err != nil {
				if err == nil || (err.Error() != tC.err.Error()) {
					t.Fatalf("ReadAll() error mismatch:\ngot %v\nwant %v", err, tC.err)
				}
				if out != nil {
					t.Fatalf("ReadAll() output:\ngot %v\nwant nil", out)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected ReadAll() error: %v", err)
				}
				if !reflect.DeepEqual(out, tC.out) {
					t.Fatalf("ReadAll() output:\ngot %v\nwant %v", out, tC.out)
				}
			}

			r = newReader(tC)
			for recNum := 0; ; recNum++ {
				rec, err := r.Read()

				var wantErr error
				if tC.err != nil {
					wantErr = tC.err
				} else if recNum >= len(tC.out) {
					wantErr = io.EOF
				}

				if err != nil && err.Error() != wantErr.Error() {
					t.Fatalf("Read() error at record %d:\ngot %v\nwant %v", recNum, err, wantErr)
				}
				if err != nil {
					break
				}
				if got, want := rec, tC.out[recNum]; !reflect.DeepEqual(got, want) {
					t.Errorf("Read vs ReadAll mismatch:\ngot %v\nwant %v", got, want)
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	for _, tC := range readTests {
		if tC.fDelim != 0 || tC.kvDelim != 0 {
			// read string creates a reader with default delimiters
			// tests with custom delimters fail
			continue
		}
		out, err := tuples.ReadString(tC.in)
		if tC.err != nil {
			if err == nil || (err.Error() != tC.err.Error()) {
				t.Fatalf("ReadString() error mismatch:\ngot %v\nwant %v", err, tC.err)
			}
			if out != nil {
				t.Fatalf("ReadString() output:\ngot %v\nwant nil", out)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected ReadString() error: %v", err)
			}
			if !reflect.DeepEqual(out, tC.out) {
				t.Fatalf("ReadString() output:\ngot %v\nwant %v", out, tC.out)
			}
		}
	}
}
