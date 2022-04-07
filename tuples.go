package tuples

import (
	"bufio"
	"io"
	"strings"
)

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
