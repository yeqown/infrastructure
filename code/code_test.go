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
		{
			name: "case1",
			args: args{code: CodeOk, message: "成功"},
			want: &CodeInfo{CodeOk, "成功"},
		},
		{
			name: "case2",
			args: args{code: CodeOk, message: ""},
			want: &CodeInfo{CodeOk, messages[CodeOk]},
		},
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
		{
			name: "case1",
			args: args{code: CodeOk},
			want: &CodeInfo{CodeOk, messages[CodeOk]},
		},
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
		{
			name: "case1",
			args: args{code: CodeOk},
			want: messages[CodeOk],
		},
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
	type testStruct struct {
		CodeInfo
		otherFiled string
	}

	// type testStruct2 struct {
	// 	CodeInfo   *CodeInfo
	// 	otherFiled string
	// }

	type args struct {
		v  interface{}
		ci *CodeInfo
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// {
		// 	name: "struct case should be panic",
		// 	args: args{
		// 		v:  testStruct{},
		// 		ci: GetCodeInfo(CodeOk),
		// 	},
		// 	want: testStruct{CodeInfo: *(GetCodeInfo(CodeOk))},
		// },
		{
			name: "struct ptr case",
			args: args{
				v:  &testStruct{},
				ci: GetCodeInfo(CodeOk),
			},
			want: &testStruct{CodeInfo: *(GetCodeInfo(CodeOk))},
		},
		// {
		// 	name: "struct ptr and codeinfo ptr case",
		// 	args: args{
		// 		v:  &testStruct2{},
		// 		ci: GetCodeInfo(CodeOk),
		// 	},
		// 	want: &testStruct2{CodeInfo: GetCodeInfo(CodeOk)},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FillCodeInfo(tt.args.v, tt.args.ci); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillCodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
