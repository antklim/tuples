package tuples

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ScannerError describes an error that occurred while scanning a tuple.
type ScannerError interface {
	error
	ScanFailed() bool
}

type scannerError struct {
	error
}

func (e scannerError) ScanFailed() bool { return true }

const (
	idxKey = iota
	idxVal
)

const (
	defaultFieldsDelimiter   = ','
	defaultKeyValueDelimiter = '='
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

var defaultScannerOptions = scannerOptions{fd: ',', kvd: '='}

type scanner struct {
	s     *bufio.Scanner
	state int
	pos   int
	err   error
	opts  scannerOptions
}

// TODO(feat): add options
func newScanner(r io.Reader, opts ...scannerOpt) *scanner {
	sopts := defaultScannerOptions
	for _, opt := range opts {
		opt(&sopts)
	}

	bufscan := bufio.NewScanner(r)
	bufscan.Split(bufio.ScanWords)
	s := &scanner{
		s:    bufscan,
		opts: sopts,
	}

	return s
}

// next moves the scanner along the tuples values. It returns false if scanning
// finished or error occurred. Call tupple() to get scanned values.
// For example:
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
			s.err = scannerError{err}
		}
	}
	s.pos++
	return s.state != scanDone
}

func (s *scanner) nextTimes(n int) bool {
	for ; s.next() && n > 1; n-- {
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
			s.err = scannerError{fmt.Errorf("tuples: tuple #%d invalid field #%d", s.pos, i+1)}
			return nil, s.err
		}
		tuple = append(tuple, []string{kv[idxKey], kv[idxVal]})
	}
	return tuple, nil
}

func splitFunc(dlm rune) func(rune) bool {
	return func(r rune) bool { return r == dlm }
}

type scannerOpt func(*scannerOptions)

func withFieldsDelimiter(d rune) scannerOpt {
	return func(so *scannerOptions) { so.fd = d }
}

func withKeyValueDelimiter(d rune) scannerOpt {
	return func(so *scannerOptions) { so.kvd = d }
}
