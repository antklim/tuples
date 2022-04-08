package tuples

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var errInvalidDelim = errors.New("tuples: invalid fields or key values delimiter")

type Reader struct {
	FieldsDelimiter rune
	KeyValDelimiter rune
	s               *bufio.Scanner
}

func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	return &Reader{
		FieldsDelimiter: ',',
		KeyValDelimiter: '=',
		s:               s,
	}
}

func (r *Reader) Read() (tuple []string, err error) {
	return r.readTuple()
}

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
	if r.FieldsDelimiter == r.KeyValDelimiter {
		return nil, errInvalidDelim
	}

	if r.s.Scan() {
		tuple := r.readFields(r.s.Text())
		return tuple, nil
	}
	if err := r.s.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}

func (r *Reader) readFields(s string) []string {
	var fieldValues []string
	fields := strings.FieldsFunc(s, func(c rune) bool { return c == r.FieldsDelimiter })
	for _, f := range fields {
		kv := strings.FieldsFunc(f, func(c rune) bool { return c == r.KeyValDelimiter })
		fieldValues = append(fieldValues, kv[1])
	}
	return fieldValues
}

func ReadString(s string) ([][]string, error) {
	r := NewReader(strings.NewReader(s))
	return r.ReadAll()
}
