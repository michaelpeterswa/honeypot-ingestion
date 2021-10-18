package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var settings Settings
var redisConn *RedisConn
var influxConn *InfluxConn
var ipinfoConn *IPInfoConn
var ctx = context.Background()

func main() {
	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		log.Println("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		log.Println("Error unmarshalling settings")
	}

	initLogger(settings.ZapLevel)
	redisConn = initRedis()
	influxConn = initInflux()
	ipinfoConn = initIPInfo()

	for {
		logEntry := redisConn.client.BRPop(ctx, time.Minute, settings.CowrieKey)
		ProcessCowrieLogEntry(logEntry)
	}
}

func ProcessCowrieLogEntry(obj *redis.StringSliceCmd) {
	var tmp map[string]interface{}
	result, err := obj.Result()
	if err == redis.Nil {
		logger.Debug("Queue is currently empty...", zap.String("key", settings.CowrieKey))
		return
	} else if err != nil {
		logger.Error("Unable to get result from *redis.StringSliceCmd")
		return
	}
	data := result[1]
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		logger.Error("Unable to get result from *redis.StringSliceCmd")
		return
	}
	eventid := tmp["eventid"]
	logger.Debug("Event Received!", zap.Any("eventid", eventid))

	switch eventid {
	case "cowrie.login.success":
		var cls CowrieLoginSuccess
		json.Unmarshal([]byte(data), &cls)
		geo := getGeoIPInfo(cls.SrcIP)
		influxConn.writeCowrieLoginSuccess(cls, geo)
	case "cowrie.login.failed":
		var clf CowrieLoginFailed
		json.Unmarshal([]byte(data), &clf)
		geo := getGeoIPInfo(clf.SrcIP)
		influxConn.writeCowrieLoginFailed(clf, geo)
	case "cowrie.session.connect":
		var csc CowrieSessionConnect
		json.Unmarshal([]byte(data), &csc)
		geo := getGeoIPInfo(csc.SrcIP)
		influxConn.writeCowrieSessionConnect(csc, geo)
	case "cowrie.session.params":
		break
	case "cowrie.session.closed":
		break
	case "cowrie.session.file_download":
		break
	case "cowrie.command.input":
		var cci CowrieCommandInput
		json.Unmarshal([]byte(data), &cci)
		geo := getGeoIPInfo(cci.SrcIP)
		influxConn.writeCowrieCommandInput(cci, geo)
	case "cowrie.command.failed":
		break
	case "cowrie.direct-tcpip.request":
		break
	case "cowrie.direct-tcpip.data":
		break
	case "cowrie.client.fingerprint":
		break
	case "cowrie.client.kex":
		break
	case "cowrie.client.version":
		break
	case "cowrie.log.closed":
		break
	default:
		break
	}
}
