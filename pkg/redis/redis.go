package redisCache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func New(host string, port string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   db,
	})
	return rdb
}
