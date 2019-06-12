package lang

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_IsJSON(t *testing.T) {
	q := "a=1&b=2&c={'a': 'a', 'b': 'b'}"

	json := "{\"a\": \"a\", \"b\": \"b\"}"
	if b := IsJSON(q); b {
		t.Error("should be json")
	}
	if b := IsJSON(json); !b {
		t.Error("should be json")
	}
}

func TestPostForm(t *testing.T) {
	type args struct {
		URL  string
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostForm(tt.args.URL, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSON(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSON(tt.args.s); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseURLQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseURLQuery(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURLQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseJSON(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data2urlValues(t *testing.T) {
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want url.Values
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := data2urlValues(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data2urlValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
