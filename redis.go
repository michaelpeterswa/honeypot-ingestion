package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConn struct {
	client redis.Client
}

func initRedis() *RedisConn {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.RedisURL, settings.RedisPort),
		Password: settings.RedisPassword,
		DB:       0,
	})
	return &RedisConn{*rdb}
}
