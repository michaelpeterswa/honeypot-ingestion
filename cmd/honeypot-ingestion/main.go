package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/michaelpeterswa/honeypot-ingestion/internal/db"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/geo"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/kv"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/logging"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/structs"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type HoneypotIngestion struct {
	Logger     *zap.Logger
	Settings   structs.Settings
	RedisConn  *kv.RedisConn
	InfluxConn *db.InfluxConn
	IPInfoConn *geo.IPInfoConn
}

func NewHoneypotIngestion(l *zap.Logger, s structs.Settings, r *kv.RedisConn, i *db.InfluxConn, ip *geo.IPInfoConn) HoneypotIngestion {
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
		log.Println("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		log.Println("Error unmarshalling settings")
	}

	logger, err := logging.InitLogger(settings.ZapLevel)
	if err != nil {
		log.Fatalln(err)
	}

	honeypotIngestion := NewHoneypotIngestion(logger, settings, kv.InitRedis(settings), db.InitInflux(logger, settings), geo.InitIPInfo(settings))

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
		geo := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, cls.SrcIP)
		hpi.InfluxConn.WriteCowrieLoginSuccess(hpi.Logger, hpi.Settings, cls, geo)
	case "cowrie.login.failed":
		var clf structs.CowrieLoginFailed
		json.Unmarshal([]byte(data), &clf)
		geo := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, clf.SrcIP)
		hpi.InfluxConn.WriteCowrieLoginFailed(hpi.Logger, hpi.Settings, clf, geo)
	case "cowrie.session.connect":
		var csc structs.CowrieSessionConnect
		json.Unmarshal([]byte(data), &csc)
		geo := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, csc.SrcIP)
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
		geo := geo.GetGeoIPInfo(ctx, hpi.Logger, hpi.RedisConn, hpi.IPInfoConn, cci.SrcIP)
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
