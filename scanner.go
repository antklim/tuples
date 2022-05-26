package tuples

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var (
	errEqualDelimiters          = errors.New("fields and key-value delimiters are equal")
	errInvalidFieldsDelimiter   = errors.New("invalid fields delimiter")
	errInvalidKeyValueDelimiter = errors.New("invalid key-value delimiter")
)

// ScannerError describes an error that occurred while scanning a tuple.
type ScannerError struct {
	err error
}

func (e *ScannerError) Error() string {
	return fmt.Sprintf("tuples: scan failed: %s", e.err)
}

func (e *ScannerError) Unwrap() error {
	return e.err
}

// InvalidScannerOptionError describes an error that occurred while initializing
// scanner with invalid options.
type InvalidScannerOptionError struct {
	err error
}

func (e *InvalidScannerOptionError) Error() string {
	return fmt.Sprintf("tuples: invalid delimiters: %s", e.err)
}

func (e *InvalidScannerOptionError) Unwrap() error {
	return e.err
}

const (
	idxKey = iota
	idxVal
)

const (
	scanReady = iota
	scanTuple
	scanDone
)

type scannerOptions struct {
	fd  rune // fields delimiter
	kvd rune // key-values delimiter
}

func (so *scannerOptions) validate() error {
	if so.fd == so.kvd {
		return errEqualDelimiters
	}
	if !validDelim(so.fd) {
		return errInvalidFieldsDelimiter
	}
	if !validDelim(so.kvd) {
		return errInvalidKeyValueDelimiter
	}
	return nil
}

func validDelim(r rune) bool {
	return r != 0 && utf8.ValidRune(r) && r != utf8.RuneError
}

var defaultScannerOptions = scannerOptions{fd: ',', kvd: '='}

type scanner struct {
	s     *bufio.Scanner
	state int
	pos   int
	err   error
	opts  scannerOptions
}

func newScanner(r io.Reader, opts ...scannerOption) (*scanner, error) {
	sopts := defaultScannerOptions
	for _, opt := range opts {
		opt(&sopts)
	}

	if err := sopts.validate(); err != nil {
		return nil, &InvalidScannerOptionError{err}
	}

	bufscan := bufio.NewScanner(r)
	bufscan.Split(bufio.ScanWords)
	s := &scanner{
		s:    bufscan,
		opts: sopts,
	}

	return s, nil
}

// next moves the scanner along the tuples values. It returns false if scanning
// finished or error occurred. Call tupple() to get scanned values.
// next returns false only when the scanner is finished. It means that even
// when the last tuple scanned the next returns true and should be called once
// more to make sure that there is nothing left. next does not scan tuple
// proactively.
//
// Example of scanner use:
//
//	in := strings.NewReader("name=Rob,lname=Doe name=Bob,lname=Smith")
//	s := newScanner(in)
//	for s.next() {
//		fmt.Println(s.tuple())
//	}
//
//	// Output:
//	// [[name Rob] [lname Doe]]
//	// [[name Bob] [lname Smith]]
func (s *scanner) next() bool {
	if s.err != nil {
		s.state = scanDone
	}
	if s.state == scanDone {
		return false
	}
	if s.state == scanReady {
		s.state = scanTuple
	}
	if !s.s.Scan() {
		s.state = scanDone
		if err := s.s.Err(); err != nil {
			s.err = &ScannerError{err}
		}
	}
	s.pos++
	return s.state != scanDone
}

func (s *scanner) nextTimes(n int) bool {
	for ; n >= 1 && s.next(); n-- {
	}
	return s.state != scanDone
}

func (s *scanner) tuple() ([][]string, error) {
	// It splits "name=John,lname=Doe,age=17" to ["name=John", "lname=Doe", "age=17"].
	fields := strings.FieldsFunc(s.s.Text(), splitFunc(s.opts.fd))
	var tuple [][]string
	for i, f := range fields {
		// It splits "name=John" into ["name", "John"].
		kv := strings.FieldsFunc(f, splitFunc(s.opts.kvd))
		if len(kv) != 2 { // nolint: gomnd
			s.err = &ScannerError{fmt.Errorf("tuple #%d invalid field #%d", s.pos, i+1)}
			return nil, s.err
		}
		tuple = append(tuple, []string{kv[idxKey], kv[idxVal]})
	}
	return tuple, nil
}

func splitFunc(dlm rune) func(rune) bool {
	return func(r rune) bool { return r == dlm }
}

type scannerOption func(*scannerOptions)

func withFieldsDelimiter(d rune) scannerOption {
	return func(so *scannerOptions) { so.fd = d }
}

func withKeyValueDelimiter(d rune) scannerOption {
	return func(so *scannerOptions) { so.kvd = d }
}
