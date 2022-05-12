package vertigo

import (
	"errors"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
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

// extractStructField takes a reflect.Value and a valueMapper and
// returns the struct field at valueMapper.fieldIdx. If v is a reflect.Ptr,
// it will extract the value at address v.
func extractStructField(v reflect.Value, lookup valueMapper) reflect.Value {
	if isValuePointer(v) {
		v = reflect.Indirect(v)
	}
	return v.Field(lookup.fieldIdx)
}

// setStructField sets the struct field within the destination struct being scanned.
func setStructField(fv *aiplatformpb.FeatureValue, structField reflect.Value) {
	switch fv.Value.(type) {
	case *aiplatformpb.FeatureValue_BoolValue:
		bv := fv.GetBoolValue()
		if isValuePointer(structField) {
			structField.Set(reflect.ValueOf(&bv))
		} else {
			structField.SetBool(bv)
		}

	case *aiplatformpb.FeatureValue_BoolArrayValue:
		structField.Set(reflect.ValueOf(fv.GetBoolArrayValue().Values))

	case *aiplatformpb.FeatureValue_Int64Value:
		iv := fv.GetInt64Value()
		if isValuePointer(structField) {
			structField.Set(reflect.ValueOf(&iv))
		} else {
			structField.SetInt(iv)
		}

	case *aiplatformpb.FeatureValue_Int64ArrayValue:
		structField.Set(reflect.ValueOf(fv.GetInt64ArrayValue().Values))

	case *aiplatformpb.FeatureValue_DoubleValue:
		dv := fv.GetDoubleValue()
		if isValuePointer(structField) {
			structField.Set(reflect.ValueOf(&dv))
		} else {
			structField.SetFloat(dv)
		}

	case *aiplatformpb.FeatureValue_DoubleArrayValue:
		structField.Set(reflect.ValueOf(fv.GetDoubleArrayValue().Values))

	case *aiplatformpb.FeatureValue_StringValue:
		stringV := fv.GetStringValue()
		if isValuePointer(structField) {
			structField.Set(reflect.ValueOf(&stringV))
		} else {
			structField.SetString(stringV)
		}

	case *aiplatformpb.FeatureValue_StringArrayValue:
		structField.Set(reflect.ValueOf(fv.GetStringArrayValue().Values))

	case *aiplatformpb.FeatureValue_BytesValue:
		structField.SetBytes(fv.GetBytesValue())
	}
}

// loadMap loads a vertex tag into it's respective field index, field name, and type from an interface.
func loadMap(dst interface{}) map[string]valueMapper {
	provided := reflect.ValueOf(dst)
	ind := reflect.Indirect(provided)
	providedType := ind.Type()
	vm := map[string]valueMapper{}
	for i := 0; i < ind.NumField(); i++ {
		structField := providedType.Field(i)
		tv, ok := structField.Tag.Lookup(vertexTag)
		if !ok || tv == "-" {
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
