package mgo

import (
	"fmt"
	"strings"
	"time"

	"github.com/yeqown/infrastructure/types"
	mgov2 "gopkg.in/mgo.v2"
)

// ConnectMgo .
func ConnectMgo(cfg *types.MgoConfig) (*mgov2.Database, error) {
	dialInfo := mgov2.DialInfo{
		Addrs:     strings.Split(cfg.Addrs, ","),
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Database:  cfg.Database,
		Username:  cfg.Username,
		Password:  cfg.Password,
		PoolLimit: cfg.PoolLimit,
	}

	// connect db
	session, err := mgov2.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("Mongo connected, address: " + cfg.Addrs)

	// settings
	session.SetMode(mgov2.Strong, true)
	session.SetSocketTimeout(time.Duration(5 * time.Second))
	return session.DB(cfg.Database), nil
}
