package gormic

import (
	"fmt"
	"sync"

	"github.com/yeqown/infrastructure/types"

	"github.com/jinzhu/gorm"
)

// ConnectMysql build a connection to mysql
func ConnectMysql(c *types.MysqlC) (db *gorm.DB, err error) {
	if db, err = gorm.Open("mysql", c.String()); err != nil {
		return nil, err
	}
	// db setting
	db.DB().SetMaxOpenConns(c.Pool)
	db.DB().SetMaxIdleConns(c.Pool / 2)
	db.LogMode(false)
	db.SingularTable(true)

	if err = db.DB().Ping(); err != nil {
		return nil, fmt.Errorf("could not ping mysql server: %v", err)
	}

	return db, nil
}

// TableConn ... connection to table
type TableConn struct {
	*gorm.DB
	*sync.RWMutex
}

// ConnectSqlite3 .
func ConnectSqlite3(c *types.SQLite3Config) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", c.Name)
	if err != nil {
		return nil, err
	}
	if err = db.DB().Ping(); err != nil {
		return nil, fmt.Errorf("could not ping mysql server: %v", err)
	}

	return db, nil
}

type healthchecker struct {
	db *gorm.DB
}

func (hc *healthchecker) Check() types.HealthInfo {
	var info = types.NewHealthInfo()
	info.Healthy = true
	if err := hc.db.DB().Ping(); err != nil {
		info.Healthy = false
		info.Meta["error"] = err.Error()
	}
	return info
}

// NewHealthChecker .
func NewHealthChecker(db *gorm.DB) types.HealthChecker {
	return &healthchecker{db: db}
}
