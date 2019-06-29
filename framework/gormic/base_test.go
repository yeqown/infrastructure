package gormic

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yeqown/infrastructure/types"
)

func TestConnectMysql(t *testing.T) {
	mysqlC := &types.MysqlC{
		UserName:  "root",
		Password:  "ncDYbAAx4mrl",
		Addr:      "192.168.2.254",
		DBName:    "bussiness",
		Charset:   "utf8",
		ParseTime: true,
		Loc:       "Local",
		Pool:      10,
	}
	gotDb, err := ConnectMysql(mysqlC)
	if err != nil {
		t.Errorf("ConnectMysql() error = %v", err)
		return
	}
	_ = gotDb
}
