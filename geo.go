package main

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ipinfo/go/v2/ipinfo"
	"go.uber.org/zap"
)

type IPInfoConn struct {
	client ipinfo.Client
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

func initIPInfo() *IPInfoConn {
	return &IPInfoConn{*ipinfo.NewClient(nil, nil, settings.IPInfoKey)}
}

func printGeoData(geo GeoData) {
	logger.Debug("Current GeoData",
		zap.String("ip", geo.IP),
		zap.String("city", geo.City),
		zap.String("country", geo.Country),
		zap.String("region", geo.Region),
		zap.String("location", geo.Location),
		zap.Float64("latitude", geo.Latitude),
		zap.Float64("longitude", geo.Longitude))
}

func getGeoIPInfo(ip string) GeoData {
	var info *ipinfo.Core
	var geo GeoData
	value := redisConn.client.Get(ctx, ip)
	result, err := value.Result()
	if err == redis.Nil {
		logger.Debug("Cache missed...", zap.Error(err))
		info, err = ipinfoConn.client.GetIPInfo(net.ParseIP(ip))
		if err != nil {
			logger.Error("Unable to parse IP and get info...", zap.Error(err), zap.String("ip", ip))
		}
		newGeo := setGeoData(info)
		bytes, err := json.Marshal(newGeo)
		if err != nil {
			logger.Error("Couldn't marshal GeoData")
		}
		redisConn.client.Set(ctx, ip, bytes, time.Hour*168)
		logger.Debug("Cache Set...", zap.ByteString("payload", bytes))
		geo = newGeo
	} else {
		logger.Debug("Cache Hit...", zap.String("data", result))
		json.Unmarshal([]byte(result), &geo)
		printGeoData(geo)
	}

	return geo
}

func setGeoData(info *ipinfo.Core) GeoData {
	lat := 0.0
	lon := 0.0
	latLon := strings.Split(info.Location, ",")
	if len(latLon) == 2 {
		lat, _ = strconv.ParseFloat(latLon[0], 64)
		lon, _ = strconv.ParseFloat(latLon[1], 64)
	}
	return GeoData{
		IP:        info.IP.String(),
		City:      info.City,
		Country:   info.CountryName,
		Region:    info.Region,
		Location:  info.Location,
		Latitude:  lat,
		Longitude: lon,
	}
}
