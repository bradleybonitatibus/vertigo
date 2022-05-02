package vertigo

import (
	"errors"
	"reflect"
)

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
	if v.Kind() == reflect.Ptr {
		return true
	}
	return false
}
