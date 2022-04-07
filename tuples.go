package tuples

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	s *bufio.Scanner
}

func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	return &Reader{
		s: s,
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
		tuple := readFields(r.s.Text())
		return tuple, nil
	}
	if err := r.s.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}

func readFields(s string) []string {
	fs := strings.Split(s, ",")

	var fields []string
	for _, v := range fs {
		kv := strings.Split(v, "=")
		fields = append(fields, kv[1])
	}
	return fields
}
