package lang

import (
	"reflect"
	"testing"
	"time"
)

func TestCurTimeFormat(t *testing.T) {
	type args struct {
		layout string
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
			if got := CurTimeFormat(tt.args.layout); got != tt.want {
				t.Errorf("CurTimeFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	type args struct {
		layout string
		value  string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTime(tt.args.layout, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTodayDate(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTodayDate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTodayDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTimeDate(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTimeDate(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimeDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
