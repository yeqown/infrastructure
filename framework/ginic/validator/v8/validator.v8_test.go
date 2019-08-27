package v8_test

import (
	"testing"

	cusVali "github.com/yeqown/infrastructure/framework/ginic/validator"
	v8 "github.com/yeqown/infrastructure/framework/ginic/validator/v8"
	"github.com/yeqown/infrastructure/framework/gormic"
	mgolib "github.com/yeqown/infrastructure/framework/mgo"
	"github.com/yeqown/infrastructure/types"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	vali "gopkg.in/go-playground/validator.v8"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserModel struct {
	gorm.Model
	Name string
}

func (m UserModel) TableName() string {
	return "users"
}

func prepareData1(sqlID uint, mgoID bson.ObjectId) (*gorm.DB, *mgo.Database, error) {
	db1, err := gormic.ConnectSqlite3(
		&types.SQLite3Config{
			Name: "./testdata/chcker.db",
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// init data
	db1.DropTableIfExists(&UserModel{})
	db1.AutoMigrate(&UserModel{})

	_m1 := &UserModel{Model: gorm.Model{ID: uint(sqlID)}, Name: "foo"}
	if err := db1.Model(&UserModel{}).Create(_m1).Error; err != nil {
		return nil, nil, err
	}

	db2, err := mgolib.ConnectMgo(&types.MgoConfig{
		Addrs:     "localhost:27017",
		Timeout:   5,
		Database:  "test",
		Username:  "",
		Password:  "",
		PoolLimit: 20,
	})

	if err != nil {
		return nil, nil, err
	}

	_m2 := bson.M{
		"_id":  mgoID,
		"name": "foo",
	}

	db2.C("user").DropCollection()
	if err := db2.C("user").Insert(_m2); err != nil {
		return nil, nil, err
	}

	return db1, db2, nil
}

func Test_Validator_ResourceCheck(t *testing.T) {
	mgoID := bson.NewObjectId()

	var foo = struct {
		SqlID uint   `json:"sql_id" validate:"reschk=sqlUser"`
		MgoID string `json:"mgo_id" validate:"reschk=mgoUser"`
	}{
		SqlID: 1,
		MgoID: mgoID.Hex(),
	}

	// prepare data
	sqlDB, mgoDB, err := prepareData1(foo.SqlID, bson.ObjectIdHex(foo.MgoID))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// init checker
	tblName := (UserModel{}).TableName()
	collName := "user"

	// init resource checker
	v8.RegisterResChk("sqlUser", cusVali.NewMySQLChecker(sqlDB, tblName))
	v8.RegisterResChk("mgoUser", cusVali.NewMgoChecker(mgoDB, collName))

	// validate struct
	var validate = vali.New(&vali.Config{})
	validate.RegisterValidation("reschk", v8.DefaultResourceCheck)
	if err = validate.Struct(foo); err != nil {
		t.Log("validate foo got err: ", err)
		t.FailNow()
	}
}
