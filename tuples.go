package tuples

import (
	"bufio"
	"io"
)

var f int

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

func (r *Reader) Read() (tuple []string, err error) {
	switch f {
	case 0:
		tuple = []string{"John", "Doe", "2000-01-01"}
	case 1:
		tuple = []string{"Bob", "Smith", "2010-10-10"}
	default:
		err = io.EOF
	}
	f++

	return
}

func (r *Reader) ReadAll() (tuples [][]string, err error) {
	tuples = [][]string{{"John", "Doe", "2000-01-01"}, {"Bob", "Smith", "2010-10-10"}}
	return
}
