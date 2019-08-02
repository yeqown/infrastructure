package redigo

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestIterKeys(t *testing.T) {
	type args struct {
		client *redis.Client
		match  string
		count  int64
		f      IterFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IterKeys(tt.args.client, tt.args.match, tt.args.count, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("IterKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
