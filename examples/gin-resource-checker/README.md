# gin-resource-checker

Mainly used to check the existence of resources

### Usage[v8]

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	// vali "gopkg.in/go-playground/validator.v9"
	vali "gopkg.in/go-playground/validator.v8"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "github.com/yeqown/infrastructure/framework/ginic/validator/v9"
	"github.com/yeqown/infrastructure/framework/ginic/validator"
	v8 "github.com/yeqown/infrastructure/framework/ginic/validator/v8"
	"github.com/yeqown/infrastructure/framework/gormic"
	mgolib "github.com/yeqown/infrastructure/framework/mgo"
	"github.com/yeqown/infrastructure/types"
)

// ignored init function

// FooForm .
type FooForm struct {
	CountryID string `form:"country_id" binding:"required,reschk=mgoCountry"`
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

	log.Printf("inserted id is: %d, %s\n", sqlID, mgoID)
	v8.RegisterResChk("sqlUser", validator.NewMySQLChecker(sqlDB, "users"))
	v8.RegisterResChk("mgoCountry", validator.NewMgoChecker(mgoDB, "contries"))

	// [WIP: gin not support validator.v9 ...]
	// register custom validation tag
	_validate := binding.Validator.Engine().(*vali.Validate)
	_validate.RegisterValidation("reschk", v8.DefaultResourceCheck)

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

	if err := e.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

```