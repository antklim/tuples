package tuples

import (
	"reflect"
	"sync"
)

// TODO(feat): create fields cache
// TODO(feat): use fields cache in scanner (https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/encoding/json/encode.go;l=1410)

type field struct {
	name string
	typ  reflect.Type
	tag  bool
}

type structFields struct {
	fields   []field
	namedIdx map[string]int
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
				tag:  tag != "",
			}
			fields = append(fields, f)
		}
	}

	nameIdx := make(map[string]int)
	for i, f := range fields {
		nameIdx[f.name] = i
	}
	return structFields{fields, nameIdx}
}

// cachedTypeFields runs typeFields and stores the result in the cache.
func cachedTypeFields(t reflect.Type) structFields {
	if cache, ok := fieldCache.Load(t); ok {
		return cache.(structFields)
	}

	cache, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return cache.(structFields)
}
