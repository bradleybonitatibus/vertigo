package vertigo

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
)

type myStruct struct {
	BoolField      bool      `vertex:"bool_field"`
	BoolPointer    *bool     `vertex:"bool_pointer"`
	Int64Field     int64     `vertex:"int_64_field"`
	Int64Pointer   *int64    `vertex:"int_64_pointer"`
	Float64Field   float64   `vertex:"float_64_field"`
	Float64Pointer *float64  `vertex:"float_64_pointer"`
	StringField    string    `vertex:"string_field"`
	StringPointer  *string   `vertex:"string_pointer"`
	ByteSlice      []byte    `vertex:"byte_slice"`
	BoolSlice      []bool    `vertex:"bool_slice"`
	Int64Slice     []int64   `vertex:"int_64_slice"`
	Float64Slice   []float64 `vertex:"float_64_slice"`
	StringSlice    []string  `vertex:"string_slice"`
}

func TestIsStructPointer(t *testing.T) {
	type test struct {
		v   interface{}
		err error
	}

	// needs to initialize then assign due to use of pointer for test case
	var s string
	s = "hello"

	type someStruct struct {
		v string
	}

	myS := someStruct{v: s}

	tests := []test{
		{
			v:   &s,
			err: errors.New("dst must be a pointer"),
		},
		{
			v:   myS,
			err: errors.New("dst must be a struct type"),
		},
		{
			v:   &myS,
			err: nil,
		},
	}

	for _, tc := range tests {
		got := isStructPointer(tc.v)
		if errors.Is(got, tc.err) && got != nil {
			t.Errorf("Expected: %v, got %v instead", tc.err, got)
		}
	}
}

func TestIsValuePointer(t *testing.T) {
	type test struct {
		v    interface{}
		want bool
	}
	var p string
	p = "hello"
	tests := []test{
		{
			v:    p,
			want: false,
		},
		{
			v:    &p,
			want: true,
		},
	}

	for _, tc := range tests {
		v := reflect.ValueOf(tc.v)
		got := isValuePointer(v)
		if tc.want != got {
			t.Errorf("Expected :%v, got %v instead", tc.want, got)
		}
	}
}

func TestSetStructField(t *testing.T) {
	type test struct {
		name        string
		fv          *aiplatformpb.FeatureValue
		dst         reflect.Value
		structField reflect.Value
	}

	myS := myStruct{}
	s := reflect.ValueOf(&myS)
	sInd := reflect.Indirect(s)

	tests := []test{
		{
			name: "bool field",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_BoolValue{
					BoolValue: true,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("BoolField"),
		},
		{
			name: "bool pointer",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_BoolValue{
					BoolValue: true,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("BoolPointer"),
		},
		{
			name: "int64 field",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_Int64Value{
					Int64Value: 100,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Int64Field"),
		},
		{
			name: "int64 pointer",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_Int64Value{
					Int64Value: 200,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Int64Pointer"),
		},
		{
			name: "float64 field",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_DoubleValue{
					DoubleValue: 100.0,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Float64Field"),
		},
		{
			name: "float64 pointer",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_DoubleValue{
					DoubleValue: 200.0,
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Float64Pointer"),
		},
		{
			name: "string field",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_StringValue{
					StringValue: "hello",
				},
			},
			dst:         s,
			structField: sInd.FieldByName("StringField"),
		},
		{
			name: "string pointer",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_StringValue{
					StringValue: "world",
				},
			},
			dst:         s,
			structField: sInd.FieldByName("StringPointer"),
		},
		{
			name: "byte slice",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_BytesValue{
					BytesValue: []byte("hello"),
				},
			},
			dst:         s,
			structField: sInd.FieldByName("ByteSlice"),
		},
		{
			name: "bool slice",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_BoolArrayValue{
					BoolArrayValue: &aiplatformpb.BoolArray{Values: []bool{true, false}},
				},
			},
			dst:         s,
			structField: sInd.FieldByName("BoolSlice"),
		},
		{
			name: "int64 slice",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_Int64ArrayValue{
					Int64ArrayValue: &aiplatformpb.Int64Array{Values: []int64{300, 400}},
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Int64Slice"),
		},
		{
			name: "float64 slice",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_DoubleArrayValue{
					DoubleArrayValue: &aiplatformpb.DoubleArray{Values: []float64{300.0, 400.0}},
				},
			},
			dst:         s,
			structField: sInd.FieldByName("Float64Slice"),
		},
		{
			name: "string slice",
			fv: &aiplatformpb.FeatureValue{
				Value: &aiplatformpb.FeatureValue_StringArrayValue{
					StringArrayValue: &aiplatformpb.StringArray{Values: []string{"Goodnight", "sweet price"}},
				},
			},
			dst:         s,
			structField: sInd.FieldByName("StringSlice"),
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setStructField(tc.fv, tc.structField)
		})
	}
}

func TestLoadMap(t *testing.T) {
	type test struct {
		name               string
		s                  interface{}
		expectedNumberKeys int
	}

	tests := []test{
		{
			name: "no vertex tags",
			s: struct {
				SomeField      string `json:"some_field"`
				SomeOtherField int    `json:"some_other_field"`
			}{},
			expectedNumberKeys: 0,
		},
		{
			name: "multiple vertex tag",
			s: struct {
				SomeField  string  `json:"some_field" vertex:"some_field"`
				OtherField float64 `json:"other_field" vertex:"other_field"`
			}{},
			expectedNumberKeys: 2,
		},
		{
			name: "ignored vertex tag",
			s: struct {
				SomeField  string  `json:"some_field" vertex:"-"`
				OtherField float64 `json:"other_field" vertex:"other_field"`
			}{},
			expectedNumberKeys: 1,
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := loadMap(tc.s)
			fmt.Println(got)
			if len(got) != tc.expectedNumberKeys {
				t.Errorf("%v: expected %v keys, got %v instead", tc.name, tc.expectedNumberKeys, len(got))
			}
		})
	}
}

func TestExtractStructField(t *testing.T) {
	type test struct {
		s         interface{}
		fieldName string
		lookup    valueMapper
	}
	tests := []test{
		{
			s:         &myStruct{},
			fieldName: "BoolField",
			lookup: valueMapper{
				fieldIdx:    0,
				fieldName:   "BoolField",
				vertexField: "bool_field",
			},
		},
		{
			s:         myStruct{},
			fieldName: "StringPointer",
			lookup: valueMapper{
				fieldIdx:    6,
				fieldName:   "StringPointer",
				vertexField: "string_pointer",
			},
		},
	}
	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.fieldName, func(t *testing.T) {
			sV := reflect.ValueOf(tc.s)
			v := extractStructField(sV, tc.lookup)
			if !v.IsValid() {
				t.Errorf("%v: was not valid, %v", tc.fieldName, v.String())
			}
		})
	}
}
