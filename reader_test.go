package tuples_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

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

// var errInvalidDelim = errors.New("tuples: invalid fields or key values delimiter")

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
	// }, {
	// 	desc:    "CustomDelimiters",
	// 	in:      "fname:John;lname:Doe;dob:2000-01-01 fname:Bob;lname:Smith;dob:2010-10-10",
	// 	out:     [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
	// 	fDelim:  ';',
	// 	kvDelim: ':',
	// }, {
	// 	desc:    "BadDelimiters1",
	// 	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	// 	err:     errInvalidDelim,
	// 	fDelim:  ':',
	// 	kvDelim: ':',
	// }, {
	// 	desc:    "BadDelimiters2",
	// 	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	// 	err:     errInvalidDelim,
	// 	fDelim:  utf8.RuneError,
	// 	kvDelim: ':',
	// }, {
	// 	desc:    "BadDelimiters3",
	// 	in:      "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	// 	err:     errInvalidDelim,
	// 	fDelim:  ':',
	// 	kvDelim: utf8.RuneError,
}}

func newReader(rt readTest) (*tuples.Reader, error) {
	r, err := tuples.NewReader(strings.NewReader(rt.in))
	if err != nil {
		return nil, err
	}
	// TODO: use new reader options to set delimiters
	// if rt.fDelim != 0 {
	// 	r.FieldsDelimiter = rt.fDelim
	// }
	// if rt.kvDelim != 0 {
	// 	r.KeyValDelimiter = rt.kvDelim
	// }
	return r, nil
}

func TestRead(t *testing.T) {
	for tI, tC := range readTests {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := newReader(tC)
			if err != nil {
				t.Fatalf("#%d: unexpected newReader() error: %v", tI, err)
			}
			for recNum := 0; ; recNum++ {
				rec, err := r.Read()

				var wantErr error
				if tC.err != nil {
					wantErr = tC.err
				} else if recNum >= len(tC.out) {
					wantErr = io.EOF
				}

				if err != nil && err.Error() != wantErr.Error() {
					t.Fatalf("#%d: Read() error at record %d:\ngot  %v\nwant %v", tI, recNum, err, wantErr)
				}
				if err != nil {
					break
				}
				// TODO(chore): add len test
				if got, want := rec, tC.out[recNum]; !reflect.DeepEqual(got, want) {
					t.Errorf("Read vs ReadAll mismatch:\ngot  %v\nwant %v", got, want)
				}
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	for tI, tC := range readTests {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := newReader(tC)
			if err != nil {
				t.Fatalf("#%d: unexpected newReader() error: %v", tI, err)
			}

			out, err := r.ReadAll()
			if tC.err != nil {
				if err == nil || (err.Error() != tC.err.Error()) {
					t.Fatalf("#%d: ReadAll() error mismatch:\ngot  %v\nwant %v", tI, err, tC.err)
				}
				if out != nil {
					t.Fatalf("#%d: ReadAll() output:\ngot  %v\nwant nil", tI, out)
				}
			} else {
				if err != nil {
					t.Fatalf("#%d: unexpected ReadAll() error: %v", tI, err)
				}
				if !reflect.DeepEqual(out, tC.out) {
					t.Fatalf("#%d: ReadAll() output:\ngot  %v\nwant %v", tI, out, tC.out)
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	for tI, tC := range readTests {
		if tC.fDelim != 0 || tC.kvDelim != 0 {
			// read string creates a reader with default delimiters
			// tests with custom delimters fail
			continue
		}
		out, err := tuples.ReadString(tC.in)
		if tC.err != nil {
			if err == nil || (err.Error() != tC.err.Error()) {
				t.Fatalf("#%d: ReadString() error mismatch:\ngot  %v\nwant %v", tI, err, tC.err)
			}
			if out != nil {
				t.Fatalf("#%d: ReadString() output:\ngot  %v\nwant nil", tI, out)
			}
		} else {
			if err != nil {
				t.Fatalf("#%d: unexpected ReadString() error: %v", tI, err)
			}
			if !reflect.DeepEqual(out, tC.out) {
				t.Fatalf("#%d: ReadString() output:\ngot  %v\nwant %v", tI, out, tC.out)
			}
		}
	}
}
