package vertigo

import (
	"errors"
	"reflect"
)

const vertexTag = "vertex"

// valueMapper is used to map struct field names to their field index, tag name, and type.
type valueMapper struct {
	fieldIdx    int
	fieldName   string
	vertexField string
	t           reflect.Type
}

// isStructPointer checks that interface x is a pointer to a struct type.
func isStructPointer(dst interface{}) error {
	t := reflect.TypeOf(dst)

	if t.Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer")
	}
	v := reflect.ValueOf(dst)
	unwrapped := reflect.Indirect(v)
	if unwrapped.Kind() != reflect.Struct {
		return errors.New("dst must be a struct type")
	}
	return nil
}

// isValuePointer returns true when a reflected Value is of Kind Ptr.
func isValuePointer(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

// setSlice is a helper function to set
func setSlice(v reflect.Value, x interface{}) {
	v.Set(reflect.MakeSlice(reflect.TypeOf(x), v.Len(), v.Len()))
}

// loadMap loads a vertex tag into it's respective field index, field name, and type from an interface.
func loadMap(dst interface{}) map[string]valueMapper {
	provided := reflect.ValueOf(dst)
	if err := isStructPointer(provided); err != nil {
		panic("provided interface was not a struct pointer")
	}
	ind := reflect.Indirect(provided)
	providedType := ind.Type()
	vm := map[string]valueMapper{}
	for i := 0; i < ind.NumField(); i++ {
		structField := providedType.Field(i)
		tv, ok := structField.Tag.Lookup(vertexTag)
		if !ok {
			continue
		}
		vm[tv] = valueMapper{
			fieldIdx:    i,
			t:           structField.Type,
			fieldName:   structField.Name,
			vertexField: tv,
		}
	}
	return vm
}
