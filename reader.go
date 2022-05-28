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

// NewReader ...
func NewReader(r io.Reader, opts ...ReaderOption) (*Reader, error) {
	ropts := defaultReaderOptions
	for _, opt := range opts {
		opt(&ropts)
	}

	s, err := newScanner(r,
		withFieldsDelimiter(ropts.fieldsDelimiter),
		withKeyValueDelimiter(ropts.keyValDelimiter))
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

// ReadString ...
func ReadString(s string) ([][]string, error) {
	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	return r.ReadAll()
}

type readerOptions struct {
	fieldsDelimiter rune
	keyValDelimiter rune
}

type ReaderOption func(*readerOptions)

func WithFieldsDelimiter(d rune) ReaderOption {
	return func(ro *readerOptions) { ro.fieldsDelimiter = d }
}

func WithKeyValueDelimiter(d rune) ReaderOption {
	return func(ro *readerOptions) { ro.keyValDelimiter = d }
}

var defaultReaderOptions = readerOptions{fieldsDelimiter: ',', keyValDelimiter: '='}
