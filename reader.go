package tuples

import (
	"io"
	"strings"
)

// TODO(chore): add comments

// Reader ...
type Reader struct {
	s *scanner
}

// TODO(chore): add test for scanner init errors
// TODO(chore): add reader options - field and key-value delimiters

// NewReader ...
func NewReader(r io.Reader) (*Reader, error) {
	s, err := newScanner(r)
	if err != nil {
		return nil, err
	}
	return &Reader{s}, nil
}

// Read ...
func (r *Reader) Read() (tuple []string, err error) {
	return r.readTuple()
}

// ReadAll ...
func (r *Reader) ReadAll() (tuples [][]string, err error) {
	for {
		tuple, err := r.readTuple()
		if err == io.EOF {
			return tuples, nil
		}
		if err != nil {
			return nil, err
		}
		tuples = append(tuples, tuple)
	}
}

func (r *Reader) readTuple() ([]string, error) {
	if r.s.next() {
		tuple, err := r.s.tuple()
		if err != nil {
			return nil, err
		}
		var fieldValues []string
		for _, field := range tuple {
			fieldValues = append(fieldValues, field[idxVal])
		}
		return fieldValues, nil
	}

	err := r.s.err
	if err == nil {
		err = io.EOF
	}
	return nil, err
}

func ReadString(s string) ([][]string, error) {
	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	return r.ReadAll()
}
