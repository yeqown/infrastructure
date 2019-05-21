package gormic

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func TestConnectMysql(t *testing.T) {
	type args struct {
		address   string
		loc       string
		parseTime string
		charset   string
		pool      int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				address:   "credit:123456@tcp(localhost:3306)/gogo_user",
				loc:       "Local",
				parseTime: "true",
				charset:   "utf8",
				pool:      10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConnectMysql(tt.args.address, tt.args.loc, tt.args.parseTime, tt.args.charset, tt.args.pool)
		})
	}
}

func TestGetDB(t *testing.T) {
	tests := []struct {
		name string
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
