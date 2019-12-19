package alg

import (
	"testing"
)

func TestStdRanker_truncFloat64(t *testing.T) {
	v := NewRanker([]float64{1231.12312, 1231.12312, 981.1230, 981.1230, 981.1230, 0.12398, 0.12398, 981.1230, 981.1230}, 3)
	r := v.(*StdRanker)

	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "case 0",
			args: args{f: 1231.123123},
			want: 1231.123,
		},
		{
			name: "case 1",
			args: args{f: 1231.123122},
			want: 1231.123,
		},
		{
			name: "case 0",
			args: args{f: 0.129},
			want: 0.129,
		},
		{
			name: "case 1",
			args: args{f: 0.1236},
			want: 0.124,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.truncFloat64(tt.args.f); got != tt.want {
				t.Errorf("StdRanker.truncFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdRanker_Rank(t *testing.T) {
	r := NewRanker(
		[]float64{1231.12312, 1231.12312, 981.1230, 981.1230, 981.1230, 0.12398, 0.12398, 981.1230, 981.1230}, // data
		3, // truncFloat64
	)

	type args struct {
		score float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case 0",
			args: args{score: 0.1238},
			want: 8,
		},
		{
			name: "case 1",
			args: args{score: 1231.12312},
			want: 1,
		},
		{
			name: "case 2",
			args: args{score: 99999.9999},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Rank(tt.args.score); got != tt.want {
				t.Errorf("StdRanker.Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdRanker_Rank2(t *testing.T) {
	r := NewRanker(
		[]float64{1231.12312, 1231.12312, 981.1230, 0.12398, 0.12398, 981.123}, // data
		3, // truncFloat64
	)

	type args struct {
		score float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case 0",
			args: args{score: 0.1238},
			want: 5,
		},
		{
			name: "case 1",
			args: args{score: 1231.12312},
			want: 1,
		},
		{
			name: "case 2",
			args: args{score: 99999.9999},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Rank(tt.args.score); got != tt.want {
				t.Errorf("StdRanker.Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}
