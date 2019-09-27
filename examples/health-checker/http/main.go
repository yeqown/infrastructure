package main

import (
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/yeqown/infrastructure/framework/gormic"
	"github.com/yeqown/infrastructure/framework/mgo"
	"github.com/yeqown/infrastructure/framework/redigo"
	"github.com/yeqown/infrastructure/types"
)

func main() {
	healthMgr := types.NewHealthMgr()
	sqliteDB, err := gormic.ConnectSqlite3(&types.SQLite3Config{
		Name: "../testdata/sqlite3.db",
	})
	if err != nil {
		panic(err)
	}
	mgoDB, err := mgo.ConnectMgo(&types.MgoConfig{
		Addrs:     "localhost:27017",
		Timeout:   5,
		Database:  "test",
		Username:  "",
		Password:  "",
		PoolLimit: 20,
	})
	if err != nil {
		panic(err)
	}
	redisC, err := redigo.ConnectRedis(&types.RedisConfig{
		Addr:     "localhost:6379",
		Password: "nopass",
		DB:       1,
	})
	if err != nil {
		panic(err)
	}

	_ = mgoDB

	healthMgr.AddChecker("sqlite", gormic.NewHealthChecker(sqliteDB), 0)
	healthMgr.AddChecker("mongo", mgo.NewHealthChecker(nil), 4)
	healthMgr.AddChecker("redis", redigo.NewHealthChecker(redisC), 0)

	http.HandleFunc("/health", healthMgr.Handler())
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
