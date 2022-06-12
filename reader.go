package tuples

import (
	"io"
	"strings"
)

// Reader describes a tuples reader.
type Reader struct {
	s *scanner
}

// NewReader creates a new instance of the Reader.
// If reader creation fails it returns error.
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

// Read reads one tuple at a time and returns fields values in the order
// they appear in the string. It returns error when read fails or when
// reached the end of the tuples input.
func (r *Reader) Read() ([]string, error) {
	return r.readTuple()
}

// ReadAll reads all tuples from the input. It returns a slice of tuples values.
// It returns error when reader initialisation failed or read process failed.
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

// ReadString reads all tuples from the string. It returns a slice of tuples
// values. It returns error when reader initialisation failed or read process
// failed.
//
// Usage:
//	tuples, err := ReadString("fname=John,lname=Doe dob=2000-01-01 age=17")
//	if err != nil {
//		return err
//	}
//	fmt.Println(tuples)
//
//	// Output:
//	// [[John Doe] [2000-01-01] [17]]
func ReadString(s string, opts ...ReaderOption) ([][]string, error) {
	r, err := NewReader(strings.NewReader(s), opts...)
	if err != nil {
		return nil, err
	}

	return r.ReadAll()
}

type readerOptions struct {
	fieldsDelimiter rune
	keyValDelimiter rune
}

// ReaderOption describes a reader option, i.e fields delimiter, key-value
// delimiter, etc.
type ReaderOption func(*readerOptions)

// WithFieldsDelimiter sets a custom fields delimiter option for reader.
// Default delimiter is ','.
func WithFieldsDelimiter(d rune) ReaderOption {
	return func(ro *readerOptions) { ro.fieldsDelimiter = d }
}

// WithFieldsDelimiter sets a custom key-value delimiter option for reader.
// Default delimiter is '='.
func WithKeyValueDelimiter(d rune) ReaderOption {
	return func(ro *readerOptions) { ro.keyValDelimiter = d }
}

var defaultReaderOptions = readerOptions{
	fieldsDelimiter: ',',
	keyValDelimiter: '=',
}
