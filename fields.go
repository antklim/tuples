package tuples

import (
	"reflect"
	"sync"
)

type field struct {
	name string
	typ  reflect.Type
	tag  string
}

type structFields struct {
	fields      []field
	fieldsByTag map[string]int
}

var fieldCache sync.Map // map[reflect.Type]structFields

// typeFields returns a list of fields that should be recognized for the given
// type.
func typeFields(t reflect.Type) structFields {
	var fields []field

	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if tag := fld.Tag.Get("tuples"); tag != "" {
			f := field{
				name: fld.Name,
				typ:  fld.Type,
				tag:  tag,
			}
			fields = append(fields, f)
		}
	}

	fieldsByTag := make(map[string]int)
	for i, f := range fields {
		fieldsByTag[f.tag] = i
	}
	return structFields{fields, fieldsByTag}
}

// cachedTypeFields runs typeFields and stores the result in the cache.
func cachedTypeFields(t reflect.Type) structFields {
	if cache, ok := fieldCache.Load(t); ok {
		return cache.(structFields)
	}

	cache, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return cache.(structFields)
}
