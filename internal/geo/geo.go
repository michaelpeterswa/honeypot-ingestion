package geo

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	zredis "676f.dev/zinc/redis"

	"github.com/go-redis/redis/v8"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/structs"
	"go.uber.org/zap"
)

type IPInfoConn struct {
	Client ipinfo.Client
}

type GeoData struct {
	IP        string
	City      string
	Country   string
	Region    string
	Location  string
	Latitude  float64
	Longitude float64
}

func InitIPInfo(settings structs.Settings) *IPInfoConn {
	return &IPInfoConn{*ipinfo.NewClient(nil, nil, settings.IPInfoKey)}
}

func PrintGeoData(logger *zap.Logger, geo *GeoData) {
	logger.Debug("Current GeoData",
		zap.String("ip", geo.IP),
		zap.String("city", geo.City),
		zap.String("country", geo.Country),
		zap.String("region", geo.Region),
		zap.String("location", geo.Location),
		zap.Float64("latitude", geo.Latitude),
		zap.Float64("longitude", geo.Longitude))
}

func GetGeoIPInfo(ctx context.Context, logger *zap.Logger, redisConn *zredis.RedisClient, ipinfoConn *IPInfoConn, ip string) (*GeoData, error) {
	var info *ipinfo.Core
	var geo *GeoData
	value, err := redisConn.Get(ctx, ip)
	if err == redis.Nil {
		logger.Debug("Cache missed...", zap.Error(err))
		info, err = ipinfoConn.Client.GetIPInfo(net.ParseIP(ip))
		if err != nil {
			return nil, err
		}
		newGeo := setGeoData(info)
		bytes, err := json.Marshal(newGeo)
		if err != nil {
			logger.Error("Couldn't marshal GeoData")
		}
		redisConn.Set(ctx, ip, string(bytes), time.Hour*168)
		logger.Debug("Cache Set...", zap.ByteString("payload", bytes))
		geo = newGeo
	} else if err != nil {
		return nil, err
	} else {
		logger.Debug("Cache Hit...", zap.String("data", value))
		err := json.Unmarshal([]byte(value), &geo)
		if err != nil {
			return nil, err
		}
		PrintGeoData(logger, geo)
	}

	return geo, nil
}

func setGeoData(info *ipinfo.Core) *GeoData {
	lat := 0.0
	lon := 0.0
	latLon := strings.Split(info.Location, ",")
	if len(latLon) == 2 {
		lat, _ = strconv.ParseFloat(latLon[0], 64)
		lon, _ = strconv.ParseFloat(latLon[1], 64)
	}
	return &GeoData{
		IP:        info.IP.String(),
		City:      info.City,
		Country:   info.CountryName,
		Region:    info.Region,
		Location:  info.Location,
		Latitude:  lat,
		Longitude: lon,
	}
}
