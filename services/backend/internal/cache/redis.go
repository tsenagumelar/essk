package cache

import "github.com/redis/go-redis/v9"

func NewRedis(address string, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
