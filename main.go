package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var settings Settings
var redisConn *RedisConn
var ctx = context.Background()

func main() {
	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		logger.Fatal("Error unmarshalling settings")
	}

	redisConn = initRedis()

	listLen, err := redisConn.client.LLen(ctx, settings.CowrieKey).Uint64()
	if err != nil {
		logger.Error("CowrieKey length failed to Uint64()", zap.Error(err))
	} else {
		fmt.Println(listLen)
	}
	var results map[string]interface{}
	curr, _ := redisConn.client.LRange(ctx, settings.CowrieKey, 0, -1).Result()
	for _, val := range curr {
		err = json.Unmarshal([]byte(val), &results)
		fmt.Println(results["src_ip"])
	}
}
