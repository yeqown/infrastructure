package dbs

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	gDB *gorm.DB
)

// GetDB get global DB instance
func GetDB() *gorm.DB {
	if gDB == nil {
		panic("ConnectMysql failed or haven't execute ConnectMysql")
	}
	return gDB
}

// ConnectMysql ... connect to mysql db
// "mysql": {
//     "pool": 20,
//     "charset": "utf8mb4",
//     "parseTime": "true",
//     "loc": "Local",
//     "address": "username:password@tcp(host:port)/dbname"
//   }
func ConnectMysql(address, loc, parseTime, charset string, pool int) {
	connAddr := fmt.Sprintf("%s?loc=%s&parseTime=%s&charset=%s", address, loc, parseTime, charset)
	db, err := gorm.Open("mysql", connAddr)
	if err != nil {
		log.Printf("Error! open mysql fail: %v\n", err.Error())
		return
	}

	db.DB().SetMaxOpenConns(pool)
	db.DB().SetMaxIdleConns(int(pool / 2))
	db.LogMode(false)
	db.SingularTable(true)

	if err = db.DB().Ping(); err != nil {
		log.Printf("Error! ping mysql fail: %v\n", err.Error())
		return
	}

	gDB = db
}
