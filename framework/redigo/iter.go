package redigo

import "github.com/go-redis/redis"

// IterFunc .
type IterFunc func(key string)

// IterKeys . 对于某一类keys迭代遍历
func IterKeys(client *redis.Client, match string, count int64, f IterFunc) error {
	// 第一次执行
	keys, cursor, err := client.Scan(0, match, count).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		f(key)
	}

	// 直到游标再次等于0
	for cursor != 0 {
		if keys, cursor, err = client.Scan(cursor, match, count).
			Result(); err != nil {
			return err
		}
		if err != nil {
			return err
		}
		for _, key := range keys {
			f(key)
		}
	}
	return nil
}
