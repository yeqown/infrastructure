// Package code to define some code
package code

import (
	"reflect"
	"testing"
)

func TestNewCodeInfo(t *testing.T) {
	type args struct {
		code    int
		message string
	}
	tests := []struct {
		name string
		args args
		want *CodeInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCodeInfo(tt.args.code, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCodeInfo(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want *CodeInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCodeInfo(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMessage(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMessage(tt.args.code); got != tt.want {
				t.Errorf("GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillCodeInfo(t *testing.T) {
	type args struct {
		v  interface{}
		ci *CodeInfo
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FillCodeInfo(tt.args.v, tt.args.ci); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillCodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
