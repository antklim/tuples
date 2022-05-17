package tuples

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// const (
// 	idxKey = iota
// 	idxVal
// )

const (
	defaultFieldsDelimiter   = ','
	defaultKeyValueDelimiter = '='
)

const (
	scanReady = iota
	scanTuple
	scanDone
)

type scanner struct {
	s     *bufio.Scanner
	state int
	pos   int
	err   error
	fd    rune // fields delimiter
	kvd   rune // key-values delimiter
}

// TODO: add options
func newScanner(r io.Reader) *scanner {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	return &scanner{
		s:   s,
		fd:  defaultFieldsDelimiter,
		kvd: defaultKeyValueDelimiter,
	}
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
		s.err = s.s.Err()
	} else {
		s.pos++
	}
	return s.state != scanDone
}

func (s *scanner) nextTimes(n int) bool {
	for ; s.next() && n > 1; n-- {
	}
	return s.state != scanDone
}

func (s *scanner) tuple() [][]string {
	// It splits "name=John,lname=Doe,age=17" to ["name=John", "lname=Doe", "age=17"].
	fields := strings.FieldsFunc(s.s.Text(), splitFunc(s.fd))
	var tuple [][]string
	for i, f := range fields {
		// It splits "name=John" into ["name", "John"].
		kv := strings.FieldsFunc(f, splitFunc(s.kvd))
		if len(kv) != 2 { // nolint: gomnd
			s.err = fmt.Errorf("tuples: tuple #%d invalid field #%d", s.pos, i+1)
			return nil
		}
		tuple = append(tuple, kv)
	}
	return tuple
}

func splitFunc(dlm rune) func(rune) bool {
	return func(r rune) bool { return r == dlm }
}
