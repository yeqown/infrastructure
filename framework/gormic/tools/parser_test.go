package tools

import (
	"testing"
)

func Test_loadGoFile(t *testing.T) {
	type args struct {
		dir      string
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
			name: "case 1",
			args: args{
				dir:      "/Users/yeqiang/go/src/github.com/yeqown/infrastructure/framework/gormic/tools/testdata",
				path:     "github.com/yeqown/infrastructure/framework/gormic/tools/testdata",
				filename: "type_model.go",
			},
			wantErr: false,
		},
	}

	isDebug = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ises, err := loadGoFiles(tt.args.dir, tt.args.path, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("loadGoFile() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				for _, is := range ises {
					t.Log(is.name, is.pkgName, is.content)
					for _, fld := range is.fields {
						t.Log(*fld)
					}
				}
			}
		})
	}
}

// func Test_cleanTag(t *testing.T) {
// 	type args struct {
// 		tag string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := cleanTag(tt.args.tag); got != tt.want {
// 				t.Errorf("cleanTag() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_parseField(t *testing.T) {
// 	type args struct {
// 		s string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *field
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := parseField(tt.args.s); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("parseField() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
