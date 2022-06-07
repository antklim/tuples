package tuples

import "errors"

// Marshal returns tuples encoding of v.
//
// Marshal supports:
//	- struct with tags
//	- pointer to a struct with tags
//	- slice of structs with tags
//	- pointer to slice of structs with tags
//	- map
//	- pointer to a map
//	- slice of maps
//	- pointer to a slice of maps
//
// In case of map, map key value used as tuple key. All map entries marshaled.
//
// Only basic types supported as values, i.e string, int, float, boolean.
// MarshalError returned in case, when unsupported type found.
func Marshal(v any) ([]byte, error) {
	return nil, errors.New("not implemented")
}
