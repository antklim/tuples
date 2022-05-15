package tuples

import (
	"bufio"
	"io"
	"strings"
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

type scanner struct {
	s     *bufio.Scanner
	state int
	err   error
	fd    rune
	kvd   rune
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

// next reads a raw tuple string and returns a slice of tuple's fields. In the
// second output parameter it returns if scan finished. For example:
//
//	s := newScanner(strings.NewReader("name=John,lname=Doe,age=17"))
//	fmt.Println(s.next()) // [[name John] [lname Doe] [age 17]], true
//
//	s := newScanner(strings.NewReader("fname=John,lname=Doe,dob=2000-01-01 fname=Bob,lname=Smith,dob=2010-10-10"))
//	fmt.Println(s.next()) // [[name John] [lname Doe] [age 17]], false
func (s *scanner) next() ([][]string, bool) {
	if s.state == scanDone {
		return nil, true
	}

	if s.state == scanReady {
		if !s.s.Scan() {
			s.state = scanDone
			s.err = s.s.Err()
			return nil, true
		}
		s.state = scanTuple
	}

	tuple := s.tuple()
	if !s.s.Scan() {
		s.state = scanDone
		s.err = s.s.Err()
	}
	return tuple, s.state == scanDone
}

func (s *scanner) tuple() [][]string { // TODO: can return [][2]string
	// It splits "name=John,lname=Doe,age=17" to ["name=John", "lname=Doe", "age=17"].
	fields := strings.FieldsFunc(s.s.Text(), splitFunc(s.fd))
	var tuple [][]string
	for _, f := range fields {
		// It splits "name=John" into ["name", "John"].
		kv := strings.FieldsFunc(f, splitFunc(s.kvd))
		tuple = append(tuple, kv)
	}
	return tuple
}

func splitFunc(dlm rune) func(rune) bool {
	return func(r rune) bool { return r == dlm }
}
