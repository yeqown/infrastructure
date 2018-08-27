package utils

import "testing"

func TestParseFloat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFloat(tt.args.s); got != tt.want {
				t.Errorf("ParseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decimal(tt.args.val); got != tt.want {
				t.Errorf("Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}
