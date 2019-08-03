package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	vali "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/yeqown/infrastructure/framework/ginic/validator.v9"
	"github.com/yeqown/infrastructure/framework/gormic"
	mgolib "github.com/yeqown/infrastructure/framework/mgo"
	"github.com/yeqown/infrastructure/types"
)

// UserModel to bind gorm.DB
type UserModel struct {
	gorm.Model
	Name string
}

// TableName .
func (UserModel) TableName() string {
	return "users"
}

func initResourceChecker(sqlID uint, mgoID bson.ObjectId) (*gorm.DB, *mgo.Database, error) {
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

	db2.C("contries").DropCollection()
	if err := db2.C("contries").Insert(_m2); err != nil {
		return nil, nil, err
	}

	return db1, db2, nil
}

// FooForm .
type FooForm struct {
	CountryID string `form:"mgo_id" binding:"required,reschk=mgoCountry"`
	UserID    uint   `form:"user_id" binding:"required,reschk=sqlUser"`
}

func main() {
	// prepare data
	var sqlID uint = 1
	var mgoID = bson.NewObjectId()
	sqlDB, mgoDB, err := initResourceChecker(sqlID, mgoID)
	if err != nil {
		panic(err)
	}
	validator.RegisterResChk("sqlUser", validator.NewMySQLChecker(sqlDB, "users"))
	validator.RegisterResChk("mgoCountry", validator.NewMgoChecker(mgoDB, "contries"))

	// [WIP: gin not support validator.v9 ...]
	// register custom validation tag
	_validate := binding.Validator.Engine().(*vali.Validate)
	_validate.RegisterValidation("reschk", validator.DefaultResourceCheck)

	e := gin.Default()
	e.GET("/resource/related_to", func(c *gin.Context) {
		var form = new(FooForm)

		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
		return
	})
}
