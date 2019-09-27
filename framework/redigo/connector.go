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
