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

type newReaderTest struct {
	desc    string
	fDelim  rune
	kvDelim rune
	err     error
}

var newReaderTests = []newReaderTest{{
	desc:    "Fails to create a reader when delimiters are the same",
	fDelim:  ':',
	kvDelim: ':',
	err:     errors.New("tuples: invalid delimiters: fields and key-value delimiters are equal"),
}, {
	desc:   "Fails to create a reader when fields delimiter is not valid",
	fDelim: utf8.RuneError,
	err:    errors.New("tuples: invalid delimiters: invalid fields delimiter"),
}, {
	desc:    "Fails to create a reader when key-value delimiter is not valid",
	kvDelim: utf8.RuneError,
	err:     errors.New("tuples: invalid delimiters: invalid key-value delimiter"),
}}

func TestNewReader(t *testing.T) {
	for tI, tC := range newReaderTests {
		t.Run(tC.desc, func(t *testing.T) {
			var opts []tuples.ReaderOption
			if tC.fDelim != 0 {
				opts = append(opts, tuples.WithFieldsDelimiter(tC.fDelim))
			}
			if tC.kvDelim != 0 {
				opts = append(opts, tuples.WithKeyValueDelimiter(tC.kvDelim))
			}

			r, err := tuples.NewReader(strings.NewReader(""), opts...)
			if err == nil || (err.Error() != tC.err.Error()) {
				t.Fatalf("#%d: NewReader error mismatch:\ngot  %v,\nwant %v", tI, err, tC.err)
			}
			var e *tuples.InvalidScannerOptionError
			if !errors.As(err, &e) {
				t.Errorf("#%d: NewReader() error is not a InvalidScannerOptionError", tI)
			}
			if r != nil {
				t.Errorf("#%d: NewReader() output:\ngot  %v\nwant nil", tI, r)
			}
		})
	}
}

type readTest struct {
	desc string
	in   string
	out  [][]string
	err  error

	fDelim  rune
	kvDelim rune
}

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
}}

func newReader(rt readTest) (*tuples.Reader, error) {
	var opts []tuples.ReaderOption
	if rt.fDelim != 0 {
		opts = append(opts, tuples.WithFieldsDelimiter(rt.fDelim))
	}
	if rt.kvDelim != 0 {
		opts = append(opts, tuples.WithKeyValueDelimiter(rt.kvDelim))
	}

	r, err := tuples.NewReader(strings.NewReader(rt.in), opts...)
	if err != nil {
		return nil, err
	}
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
				if !reflect.DeepEqual(rec, tC.out[recNum]) {
					t.Errorf("#%d: Read vs ReadAll mismatch:\ngot  %v\nwant %v", tI, rec, tC.out[recNum])
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
					t.Errorf("#%d: ReadAll() output:\ngot  %v\nwant nil", tI, out)
				}
			} else {
				if err != nil {
					t.Fatalf("#%d: unexpected ReadAll() error: %v", tI, err)
				}
				if !reflect.DeepEqual(out, tC.out) {
					t.Errorf("#%d: ReadAll() output:\ngot  %v\nwant %v", tI, out, tC.out)
				}
			}
		})
	}
}

// TODO (chore): read string error test
// TODO (chore): read error test
// TODO (chore): readAll read error test

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
				t.Errorf("#%d: ReadString() output:\ngot  %v\nwant nil", tI, out)
			}
		} else {
			if err != nil {
				t.Fatalf("#%d: unexpected ReadString() error: %v", tI, err)
			}
			if !reflect.DeepEqual(out, tC.out) {
				t.Errorf("#%d: ReadString() output:\ngot  %v\nwant %v", tI, out, tC.out)
			}
		}
	}
}
