package gormic

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/yeqown/infrastructure/types"
)

func Test_NewID(t *testing.T) {
	id := newID()
	t.Log(id, len(id))
	id = newID()
	t.Log(id, len(id))
}

func Benchmark_newID(b *testing.B) {
	// c := make(map[string]bool)
	for i := 0; i < b.N; i++ {
		_ = newID()
	}
}

/*
pkg: github.com/yeqown/server-common/framework/gormic
Benchmark_newID-8   	 1000000	      2117 ns/op	     272 B/op	       6 allocs/op

pkg: github.com/yeqown/server-common/framework/gormic
Benchmark_newID-8   	 1000000	      1966 ns/op	     144 B/op	       2 allocs/op
*/

type DemoModel struct {
	Model
	OtherColumnField uint `gorm:"other_column_field"`
}

func GetDemoModelConn(db *gorm.DB) *TableConn {
	return &TableConn{DB: db}
}

func TestModel(t *testing.T) {
	mysqlC := &types.MysqlC{
		UserName:  "root",
		Password:  "ncDYbAAx4mrl",
		Addr:      "127.0.0.1",
		DBName:    "test",
		Charset:   "utf8",
		ParseTime: true,
		Loc:       "Local",
		Pool:      10,
	}
	db, err := ConnectMysql(mysqlC)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&DemoModel{}).Error; err != nil {
		t.Fatal(err)
	}
	// create an empty model, shoud has id
	model := &DemoModel{}
	if err := GetDemoModelConn(db).Create(&model).Error; err != nil {
		t.Fatal(err)
	}

	t.Log("got model result:", *model)
	if model.ID == "" {
		t.Errorf("want id is not empty, but got: %s", model.ID)
	}
}
