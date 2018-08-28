package tools

import (
	"reflect"
	"testing"
)

func Test_loadGoFile(t *testing.T) {
	type args struct {
		filename string
		path     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "case 1",
			args:    args{path: "/Users/yeqiang/go/src/github.com/yeqown/server-common/dbs/tools/testdata", filename: "type_model.go"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadGoFiles(tt.args.path, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("loadGoFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cleanTag(t *testing.T) {
	type args struct {
		tag string
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
			if got := cleanTag(tt.args.tag); got != tt.want {
				t.Errorf("cleanTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseField(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *field
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseField(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseField() = %v, want %v", got, tt.want)
			}
		})
	}
}
