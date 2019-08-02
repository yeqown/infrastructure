package validator_test

import (
	"testing"

	"github.com/yeqown/infrastructure/framework/ginic/validator.v9"
	"github.com/yeqown/infrastructure/framework/gormic"
	"github.com/yeqown/infrastructure/framework/mgo"
	"github.com/yeqown/infrastructure/types"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/mgo.v2/bson"
)

type UserModel struct {
	gorm.Model
	Name string
}

func (m UserModel) TableName() string {
	return "users"
}

func Test_MySQLChecker(t *testing.T) {
	db, err := gormic.ConnectSqlite3(
		&types.SQLite3Config{
			Name: "./testdata/chcker.db",
		},
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// init data
	db.DropTableIfExists(&UserModel{})
	db.AutoMigrate(&UserModel{})

	// init checker
	tblName := (UserModel{}).TableName()
	checker := validator.NewMySQLChecker(db, tblName)

	// test before data exist
	if err := checker.CheckInt64(1); err == nil {
		t.Error("want err, got nil")
		t.FailNow()
	}

	// create one record
	_m := &UserModel{Name: "foo"}
	if err := db.Model(&UserModel{}).Create(_m).Error; err != nil {
		t.Error(err)
		t.FailNow()
	}

	// test after data exist
	if err := checker.CheckInt64(1); err != nil {
		t.Error("should be no err, got err: ", err)
		t.FailNow()
	}
}

func Test_MgoChecker(t *testing.T) {
	db, err := mgo.ConnectMgo(&types.MgoConfig{
		Addrs:     "localhost:27017",
		Timeout:   5,
		Database:  "test",
		Username:  "",
		Password:  "",
		PoolLimit: 20,
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	checker := validator.NewMgoChecker(db, "user")

	bid := bson.NewObjectId()
	// test before data exist
	if err := checker.Check(bid.Hex()); err == nil {
		t.Error("want err, got nil")
		t.FailNow()
	}

	_m := bson.M{
		"_id":  bid,
		"name": "foo",
	}
	if err := db.C("user").Insert(_m); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// test after data exist
	if err := checker.Check(bid.Hex()); err != nil {
		t.Error("should be no err, got err: ", err)
		t.FailNow()
	}
}
