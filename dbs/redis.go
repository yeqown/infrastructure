package dbs

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

// NewConnection ... new a connection from redis pool
func NewConnection() redis.Conn {
	if redisPool == nil {
		ConnectRedis(":6379", "", "0")
	}
	return redisPool.Get()
}

// ConnectRedis ... connect to redis
func ConnectRedis(addr, password, db string) {
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		IdleTimeout: 240 * time.Second,
		MaxIdle:     3,
	}
}
