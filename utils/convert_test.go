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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ConvertStructToMap(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ConvertStructToMap() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
