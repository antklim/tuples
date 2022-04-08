package tuples_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/antklim/tuples"
)

type readTest struct {
	Name            string
	Input           string
	FieldsDelimiter rune
	KeyValDelimiter rune
	Output          [][]string
}

var readTests = []readTest{{
	Name:   "Simple",
	Input:  "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	Name:   "EOLTest",
	Input:  "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10\n",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	Name:   "TrailingWhitespace",
	Input:  "fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10   ",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	Name:   "MultipleWhitespaces",
	Input:  "fname=John,lname=Doe,dob=2000-01-01    fname=Bob,lname=Smith,dob=2010-10-10\n",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	Name: "TabDelimited",
	Input: "fname=John,lname=Doe,dob=2000-01-01	fname=Bob,lname=Smith,dob=2010-10-10\n",
	Output: [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}, {
	Name:            "CusomDelimiters",
	Input:           "fname:John;lname:Doe;dob:2000-01-01 fname:Bob;lname:Smith;dob:2010-10-10",
	FieldsDelimiter: ';',
	KeyValDelimiter: ':',
	Output:          [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}},
}}

func newReader(rt readTest) *tuples.Reader {
	r := tuples.NewReader(strings.NewReader(rt.Input))
	if rt.FieldsDelimiter != 0 {
		r.FieldsDelimiter = rt.FieldsDelimiter
	}
	if rt.KeyValDelimiter != 0 {
		r.KeyValDelimiter = rt.KeyValDelimiter
	}
	return r
}

func TestRead(t *testing.T) {
	for _, tC := range readTests {
		t.Run(tC.Name, func(t *testing.T) {
			r := newReader(tC)
			out, err := r.ReadAll()
			if err != nil {
				t.Fatalf("unexpected ReadAll() error: %v", err)
			}
			if !reflect.DeepEqual(out, tC.Output) {
				t.Fatalf("ReadAll() output:\ngot %v\nwant %v", out, tC.Output)
			}

			r = newReader(tC)
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

func TestReadString(t *testing.T) {
	for _, tC := range readTests {
		if tC.FieldsDelimiter != 0 || tC.KeyValDelimiter != 0 {
			// read string creates a reader with default delimiters
			// tests with custom delimters fail
			continue
		}
		out, err := tuples.ReadString(tC.Input)
		if err != nil {
			t.Fatalf("unexpected ReadString() error: %v", err)
		}
		if !reflect.DeepEqual(out, tC.Output) {
			t.Fatalf("ReadString() output:\ngot %v\nwant %v", out, tC.Output)
		}
	}
}
