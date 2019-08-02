package mgo_test

import (
	"reflect"
	"testing"

	"github.com/yeqown/infrastructure/framework/mgo"

	"github.com/yeqown/infrastructure/types"
	mgov2 "gopkg.in/mgo.v2"
)

func TestConnectMgo(t *testing.T) {
	type args struct {
		cfg *types.MgoConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *mgov2.Database
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mgo.ConnectMgo(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectMgo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectMgo() = %v, want %v", got, tt.want)
			}
		})
	}
}
