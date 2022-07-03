package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"676f.dev/zinc"
	zredis "676f.dev/zinc/redis"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/db"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/geo"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/structs"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type HoneypotIngestion struct {
	Logger     *zap.Logger
	Settings   structs.Settings
	RedisConn  *zredis.RedisClient
	InfluxConn *db.InfluxConn
	IPInfoConn *geo.IPInfoConn
}

func NewHoneypotIngestion(l *zap.Logger, s structs.Settings, r *zredis.RedisClient, i *db.InfluxConn, ip *geo.IPInfoConn) HoneypotIngestion {
	return HoneypotIngestion{
		Logger:     l,
		Settings:   s,
		RedisConn:  r,
		InfluxConn: i,
		IPInfoConn: ip,
	}
}

func main() {
	var settings structs.Settings
	ctx := context.Background()

	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		log.Fatalln("Error loading settings.yaml file: ", err)
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		log.Fatalln("Error unmarshalling settings: ", err)
	}

	logger, err := zinc.InitLogger(settings.ZapLevel)
	if err != nil {
		log.Fatalln("Error initializing logger: ", err)
	}

	influx, err := zinc.InitInfluxDB(ctx, settings.InfluxAddress, settings.InfluxToken)
	if err != nil {
		logger.Fatal("Error initializing InfluxDB", zap.Error(err))
	}

	redis := zredis.NewRedisClient(logger,
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", settings.RedisURL, settings.RedisPort),
			Password: settings.RedisPassword,
			DB:       0,
		},
		influx,
		settings.InfluxOrganization,
		settings.InfluxBucket,
	)

	honeypotIngestion := NewHoneypotIngestion(logger, settings, redis, db.InitInflux(logger, influx), geo.InitIPInfo(settings))

	for {
		logEntry := honeypotIngestion.RedisConn.Client.BRPop(ctx, time.Minute, settings.CowrieKey)
		honeypotIngestion.ProcessCowrieLogEntry(ctx, logEntry)
	}
}

func (hpi HoneypotIngestion) ProcessCowrieLogEntry(ctx context.Context, obj *redis.StringSliceCmd) {
	var tmp map[string]interface{}
	result, err := obj.Result()
	if err == redis.Nil {
		hpi.Logger.Debug("Queue is currently empty...", zap.String("key", hpi.Settings.CowrieKey))
		return
	} else if err != nil {
		hpi.Logger.Error("Unable to get result from *redis.StringSliceCmd", zap.Error(err))
		return
	}
	data := result[1]
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		hpi.Logger.Error("Unable to get result from *redis.StringSliceCmd")
		return
	}
	eventid := tmp["eventid"]
	hpi.Logger.Debug("Event Received!", zap.Any("eventid", eventid))

	switch eventid {
	case "cowrie.login.success":
		var cls structs.CowrieLoginSuccess
		json.Unmarshal([]byte(data), &cls)
		geo, err := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, cls.SrcIP)
		if err != nil {
			hpi.Logger.Error("Unable to get geoip info", zap.Error(err))
			return
		}
		hpi.InfluxConn.WriteCowrieLoginSuccess(hpi.Logger, hpi.Settings, cls, geo)
	case "cowrie.login.failed":
		var clf structs.CowrieLoginFailed
		json.Unmarshal([]byte(data), &clf)
		geo, err := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, clf.SrcIP)
		if err != nil {
			hpi.Logger.Error("Unable to get geoip info", zap.Error(err))
			return
		}
		hpi.InfluxConn.WriteCowrieLoginFailed(hpi.Logger, hpi.Settings, clf, geo)
	case "cowrie.session.connect":
		var csc structs.CowrieSessionConnect
		json.Unmarshal([]byte(data), &csc)
		geo, err := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, csc.SrcIP)
		if err != nil {
			hpi.Logger.Error("Unable to get geoip info", zap.Error(err))
			return
		}
		hpi.InfluxConn.WriteCowrieSessionConnect(hpi.Logger, hpi.Settings, csc, geo)
	case "cowrie.session.params":
		break
	case "cowrie.session.closed":
		break
	case "cowrie.session.file_download":
		break
	case "cowrie.command.input":
		var cci structs.CowrieCommandInput
		json.Unmarshal([]byte(data), &cci)
		geo, err := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, cci.SrcIP)
		if err != nil {
			hpi.Logger.Error("Unable to get geoip info", zap.Error(err))
			return
		}
		hpi.InfluxConn.WriteCowrieCommandInput(hpi.Logger, hpi.Settings, cci, geo)
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
