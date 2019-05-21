package tools

import "testing"

func Test_generateFile(t *testing.T) {
	type args struct {
		ises []*innerStruct
		cfg  *outfileCfg
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				ises: []*innerStruct{
					&innerStruct{
						fields: []*field{
							{"Name", "string", "name"},
							{"Password", "string", "password"},
							{"CreateTime", "time.Time", "create_time"},
							{"UpdateTime", "time.Time", "update_time"},
						},
						content: "empty",
						name:    "UserModel",
						pkgName: "testdata",
					},
					&innerStruct{
						fields: []*field{
							{"Name", "string", "name"},
							{"Password", "string", "password"},
							{"CreateTime", "time.Time", "create_time"},
							{"UpdateTime", "time.Time", "update_time"},
						},
						content: "empty",
						name:    "User2Model",
						pkgName: "testdata",
					},
				},
				cfg: &outfileCfg{
					exportFilename:  "testdata/testfile.go.txt",
					exportPkgName:   "testdata ",
					modelImportPath: "models",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateFile(tt.args.ises, tt.args.cfg)
		})
	}
}
