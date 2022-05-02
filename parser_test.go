package vertigo

import (
	"errors"
	"reflect"
	"testing"
)

func TestIsStructPointer(t *testing.T) {
	type test struct {
		v   interface{}
		err error
	}

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
