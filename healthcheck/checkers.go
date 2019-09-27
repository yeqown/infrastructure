// Package healthcheck provide ability to
// check third-party services are alive of offline.
// currently it's simple and base
package healthcheck

import (
	"database/sql"
	"net"

	"github.com/go-redis/redis"
	mgov2 "gopkg.in/mgo.v2"
)

type healthcheckerMgo struct {
	db *mgov2.Database
}

func (hc *healthcheckerMgo) Check() HealthInfo {
	var info = NewHealthInfo()
	info.Healthy = true

	hc.db.Session.Refresh()

	if err := hc.db.Session.Ping(); err != nil {
		info.Healthy = false
		info.Meta["error"] = err.Error()
	}
	return info
}

// NewMgoChecker .
func NewMgoChecker(db *mgov2.Database) HealthChecker {
	return &healthcheckerMgo{db: db}
}

type healthcheckerRedis struct {
	client *redis.Client
}

func (hc *healthcheckerRedis) Check() HealthInfo {
	var info = NewHealthInfo()
	info.Healthy = true
	if s, err := hc.client.Ping().Result(); err != nil {
		info.Healthy = false
		info.Meta["error"] = err.Error()
		info.Meta["s"] = s
	}
	return info
}

// NewRedisChecker .
func NewRedisChecker(client *redis.Client) HealthChecker {
	return &healthcheckerRedis{client: client}
}

type healthcheckerSQL struct {
	db *sql.DB
}

func (hc *healthcheckerSQL) Check() HealthInfo {
	var info = NewHealthInfo()
	info.Healthy = true
	if err := hc.db.Ping(); err != nil {
		info.Healthy = false
		info.Meta["error"] = err.Error()
	}
	return info
}

// NewSQLChecker .
func NewSQLChecker(db *sql.DB) HealthChecker {
	return &healthcheckerSQL{db: db}
}

type healthcheckerTCP struct {
	// addr means [schema,default@tcp]://host:port
	addr string
}

func (hc *healthcheckerTCP) Check() HealthInfo {
	var info = NewHealthInfo()
	info.Healthy = true

	conn, err := net.Dial("tcp", hc.addr)
	if err != nil {
		// true: connect enconuter an error
		info.Healthy = false
		info.Meta["error"] = err.Error()
	} else {
		// true: connect ok
		conn.Close()
	}

	return info
}

// NewTCPChecker .
func NewTCPChecker(addr string) HealthChecker {
	return &healthcheckerTCP{addr: addr}
}
