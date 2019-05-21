package gormic

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

// SetLogger to open log mode and set logger
// notice that: gorm.LogWriter is an interface with Println method,
// and if logger is nil, default output to os.Stdout
func SetLogger(db *gorm.DB, logger gorm.LogWriter) *gorm.DB {
	db.LogMode(true)

	if logger == nil {
		logger = log.New(os.Stdout, "\r\n", 0)
	}

	db.SetLogger(gorm.Logger{LogWriter: logger})

	return db
}
