package redigo

import (
	"reflect"
	"testing"

	"github.com/go-redis/redis"
	"github.com/yeqown/infrastructure/types"
)

func TestConnectRedis(t *testing.T) {
	type args struct {
		cfg *types.RedisConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *redis.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnectRedis(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}
