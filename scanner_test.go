package tuples

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
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
	for tI, tC := range scanTests {
		t.Run(tC.desc, func(t *testing.T) {
			s, err := newScanner(strings.NewReader(tC.in))
			if err != nil {
				t.Fatalf("#%d: unexpected newScanner() error: %v", tI, err)
			}

			var out [][][]string
			var tuple [][]string
			// while has next
			for s.next() {
				tuple, err = s.tuple()
				if tuple != nil {
					out = append(out, tuple)
				}
			}

			if tC.err != nil {
				if err == nil || (err.Error() != tC.err.Error()) {
					t.Errorf("#%d: scan next() error mismatch:\ngot  %v\nwant %v", tI, err, tC.err)
				}
				var e ScannerError
				if !errors.As(err, &e) {
					t.Errorf("#%d: scan next() error is not a ScanError", tI)
				}
			} else {
				for i, tout := range tC.out {
					if !reflect.DeepEqual(out[i], tout) {
						t.Errorf("#%d: scan next() output #%d:\ngot  %v\nwant %v", tI, i, out[i], tout)
					}
				}
				if len(out) != len(tC.out) {
					t.Errorf("#%d: scan next() output length mismatch:\ngot  %d\nwant %d", tI, len(out), len(tC.out))
				}
			}
		})
	}
}

func TestNextAfterDone(t *testing.T) {
	s, err := newScanner(strings.NewReader("fname=John"))
	if err != nil {
		t.Fatalf("unexpected newScanner() error: %v", err)
	}

	if s.nextTimes(2) {
		t.Error("scan nextTimes(2):\ngot  true\nwant false")
	}
	out, err := s.tuple()
	if err != nil {
		t.Errorf("unexpected scan tuple() error: %v", err)
	}
	if out != nil {
		t.Errorf("scan tuple() output:\ngot  %v\nwant nil", out)
	}
}

func TestDoubleNext(t *testing.T) {
	s, err := newScanner(strings.NewReader("fname=John fname=Bob"))
	if err != nil {
		t.Fatalf("unexpected newScanner() error: %v", err)
	}

	s.nextTimes(2)

	// Tuple should return the latest data
	want := [][]string{{"fname", "Bob"}}
	out, err := s.tuple()
	if err != nil {
		t.Errorf("unexpected scan tuple() error: %v", err)
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("scan tuple() output:\ngot  %v\nwant %v", out, want)
	}
}

func TestNextTimes(t *testing.T) {
	// TODO(chore): move TestNextAfterDone and TestDoubleNext into here
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

type scannerOptTest struct {
	desc  string
	opts  []scannerOption
	sopts scannerOptions
	err   error
}

var scannerOptTests = []scannerOptTest{{
	desc:  "scanner with default settings",
	sopts: scannerOptions{fd: ',', kvd: '='},
}, {
	desc:  "scanner with custom fields delimiter",
	opts:  []scannerOption{withFieldsDelimiter(';')},
	sopts: scannerOptions{fd: ';', kvd: '='},
}, {
	desc:  "scanner with custom key-value delimiter",
	opts:  []scannerOption{withKeyValueDelimiter(':')},
	sopts: scannerOptions{fd: ',', kvd: ':'},
}, {
	desc:  "scanner with custom fields and key-value delimiters",
	opts:  []scannerOption{withFieldsDelimiter(';'), withKeyValueDelimiter(':')},
	sopts: scannerOptions{fd: ';', kvd: ':'},
}, {
	desc: "scanner with the same delimiter is not valid",
	opts: []scannerOption{withFieldsDelimiter(';'), withKeyValueDelimiter(';')},
	err:  errors.New("tuples: invalid delimiters: fields and key-value delimiters are equal"),
}, {
	desc: "scanner with invalid fields delimiter",
	opts: []scannerOption{withFieldsDelimiter(utf8.RuneError)},
	err:  errors.New("tuples: invalid delimiters: invalid fields delimiter"),
}, {
	desc: "scanner with invalid key-value delimiter",
	opts: []scannerOption{withKeyValueDelimiter(utf8.RuneError)},
	err:  errors.New("tuples: invalid delimiters: invalid key-value delimiter"),
}}

func TestScannerOptions(t *testing.T) {
	for tI, tC := range scannerOptTests {
		t.Run(tC.desc, func(t *testing.T) {
			s, err := newScanner(nil, tC.opts...)
			if tC.err != nil {
				if err == nil || (err.Error() != tC.err.Error()) {
					t.Errorf("#%d: newScanner() error mismatch:\ngot  %v\nwant %v", tI, err, tC.err)
				}
				var e *InvalidScannerOptionError
				if !errors.As(err, &e) {
					t.Errorf("#%d: newScanner() error is not a InvalidScannerOptionError", tI)
				}
				if errors.Unwrap(e) == nil {
					t.Errorf("#%d: newScanner() error should wrap original error", tI)
				}
				if s != nil {
					t.Errorf("#%d: newScanner() output:\ngot  %v\nwant nil", tI, s)
				}
			} else {
				if err != nil {
					t.Fatalf("#%d: unexpected newScanner() error: %v", tI, err)
				}
				if !reflect.DeepEqual(s.opts, tC.sopts) {
					t.Errorf("#%d: newScanner() invalid scanner options:\ngot %v\nwant %v", tI, s.opts, tC.sopts)
				}
			}
		})
	}

	// t.Run("scanner does not allow invalid UTF symbols ad delimiter", func(t *testing.T) {
	// 	s := newScanner(nil, withFieldsDelimiter('invalid utf symbol'))
	// 	// should return error
	// })
}
