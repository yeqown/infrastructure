package redigo

import (
	"github.com/go-redis/redis"
	"github.com/yeqown/infrastructure/types"
)

// ConnectRedis build a connection to redis
// redis-cli -h 106.14.168.202 -p 31079 -a class100-redis-password
func ConnectRedis(cfg *types.RedisConfig) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:     cfg.Addr,     // use default Addr
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	}
	db := redis.NewClient(opt)
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type healthchecker struct {
	client *redis.Client
}

func (hc *healthchecker) Check() types.HealthInfo {
	var info = types.NewHealthInfo()
	info.Healthy = true
	if s, err := hc.client.Ping().Result(); err != nil {
		info.Healthy = false
		info.Meta["error"] = err.Error()
		info.Meta["s"] = s
	}
	return info
}

// NewHealthChecker .
func NewHealthChecker(client *redis.Client) types.HealthChecker {
	return &healthchecker{client: client}
}
