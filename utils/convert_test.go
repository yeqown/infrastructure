package utils

import (
	"reflect"
	"testing"
)

func TestConvertStructToMap(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantOut map[string]interface{}
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				struct {
					Field1 int `json:"field_1"`
					Field2 int
				}{1, 2},
			},
			wantOut: map[string]interface{}{
				"field_1": 1,
				"Field2":  2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ConvertStructToMap(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ConvertStructToMap() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	type IStruct struct {
		Field string
		Slice []int
	}
	type Struct struct {
		Ptr    *IStruct
		IS     IStruct
		FieldS string
	}

	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 1",
			args: args{
				v: &Struct{
					Ptr:    &IStruct{"(*hj1bv232AAa", []int{123, 123123, 12312}},
					IS:     IStruct{"HJKAJSBAK", []int{12312, 123, 1231, 231231}},
					FieldS: "*&AGSIASBUAIS123123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ToLower(tt.args.v)
			t.Logf("result: %v", tt.args.v)
			iv := tt.args.v.(*Struct)
			t.Logf("result: %v", iv.Ptr)
		})
	}
}

func Test_mustbePtr(t *testing.T) {
	ptr := new(int)
	*ptr = 19
	nonPtr := 20

	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				in: ptr,
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				in: nonPtr,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustbePtr(tt.args.in); got != tt.want {
				t.Errorf("mustbePtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typeEqual(t *testing.T) {
	type args struct {
		v    reflect.Value
		kind reflect.Kind
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				v:    reflect.ValueOf("askajhsk"),
				kind: reflect.String,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeEqual(tt.args.v, tt.args.kind); got != tt.want {
				t.Errorf("typeEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
