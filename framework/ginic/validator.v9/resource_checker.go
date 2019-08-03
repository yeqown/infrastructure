package validator

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	_ ResourceChecker = MySQLChecker{}
	_ ResourceChecker = MgoChecker{}
)

// ResourceChecker .
type ResourceChecker interface {
	Check(id string) error
	CheckInt64(id int64) error
	Tag() string
}

// RegisterResChk to bind name with checker
func RegisterResChk(name string, ic ResourceChecker) {
	_checkers[name] = ic
}

// MySQLChecker is the default chcker for resource in MySQL DB
type MySQLChecker struct {
	db      *gorm.DB
	tblName string
}

// CheckInt64 of MySQLChecker .
func (c MySQLChecker) CheckInt64(id int64) error {
	cnt := 0
	err := c.db.Table(c.tblName).Where("id = ?", id).Count(&cnt).Error
	// fmt.Printf("err: %v, cnt: %d", err, cnt)
	if err == nil && cnt == 1 {
		return nil
	}
	return errors.Errorf("could not find resource with: %d", id)
}

// Check of MySQLChecker .
func (c MySQLChecker) Check(s string) error {
	return errors.New("do not support")
}

// Tag of MySQLChecker .
func (c MySQLChecker) Tag() string {
	return c.tblName
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

// CheckInt64 of MgoChecker .
func (c MgoChecker) CheckInt64(id int64) error {
	return errors.New("do not support")
}

// Tag of MgoChecker .
func (c MgoChecker) Tag() string {
	return c.collName
}

// NewMgoChecker .
func NewMgoChecker(db *mgo.Database, collName string) MgoChecker {
	return MgoChecker{db, collName}
}
