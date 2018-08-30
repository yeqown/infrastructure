package dbs

import (
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestConnectRedis(t *testing.T) {
	type args struct {
		addr     string
		password string
		db       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{
				addr:     "127.0.0.1",
				password: "",
				db:       "0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConnectRedis(tt.args.addr, tt.args.password, tt.args.db)
		})
	}
}

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name string
		want redis.Conn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConnection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
