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

type typFields struct {
	fields      []field
	fieldsByTag map[string]int
}

var fieldsCache sync.Map // map[reflect.Type]typFields

// typeFields returns a list of fields that should be recognized for the given
// type.
func typeFields(t reflect.Type) typFields {
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
	return typFields{fields, fieldsByTag}
}

// cachedTypeFields runs typeFields and stores the result in the cache.
func cachedTypeFields(t reflect.Type) typFields {
	if cache, ok := fieldsCache.Load(t); ok {
		return cache.(typFields)
	}

	cache, _ := fieldsCache.LoadOrStore(t, typeFields(t))
	return cache.(typFields)
}
