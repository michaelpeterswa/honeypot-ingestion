package kv

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/structs"
)

type RedisConn struct {
	Client redis.Client
}

func InitRedis(settings structs.Settings) *RedisConn {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.RedisURL, settings.RedisPort),
		Password: settings.RedisPassword,
		DB:       0,
	})
	return &RedisConn{*rdb}
}
