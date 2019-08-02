package validator

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _checkers map[string]ResourceChecker

func init() {
	_checkers = make(map[string]ResourceChecker)
}

// ResourceChecker .
type ResourceChecker interface {
	Check(id int64) error
}

// Register to bind name with checker
func Register(name string, ic ResourceChecker) {
	_checkers[name] = ic
}

// MySQLChecker is the default chcker for resource in MySQL DB
type MySQLChecker struct {
	db      *gorm.DB
	tblName string
}

// Check of  MySQLChecker .
func (c MySQLChecker) Check(id int64) error {
	cnt := 0
	err := c.db.Table(c.tblName).Where("id = ?", id).Count(&cnt).Error
	fmt.Printf("err: %v, cnt: %d", err, cnt)
	if err == nil && cnt == 1 {
		return nil
	}
	return errors.Errorf("could not find resource with: %d ", id)
}

// NewMySQLChecker .
func NewMySQLChecker(db *gorm.DB, tblName string) MySQLChecker {
	return MySQLChecker{db, tblName}
}

// MgoChecker is the default chcker for resource in MySQL DB
type MgoChecker struct {
	db       *mgo.Database
	collName string
}

// Check of  MgoChecker .
func (c MgoChecker) Check(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.Errorf("invalid ObjectId hex string")
	}

	cnt, err := c.db.C(c.collName).FindId(bson.ObjectIdHex(id)).Count()
	if err == nil && cnt == 1 {
		return nil
	}

	return errors.Errorf("could not find resource with: %s", id)
}

// NewMgoChecker .
func NewMgoChecker(db *mgo.Database, collName string) MgoChecker {
	return MgoChecker{db, collName}
}
