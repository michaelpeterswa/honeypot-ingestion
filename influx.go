package main

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

type InfluxConn struct {
	client influxdb2.Client
}

func initInflux() *InfluxConn {
	logger.Info("Creating InfluxDB v2 connection...")
	client := influxdb2.NewClient(settings.InfluxAddress, settings.InfluxToken)
	return &InfluxConn{client}
}

func getOrganization() string {
	return settings.InfluxOrganization
}

func getBucket() string {
	return settings.InfluxBucket
}

func (conn *InfluxConn) writeCowrieLoginSuccess(cls CowrieLoginSuccess, geo GeoData) {
	p := influxdb2.NewPointWithMeasurement(settings.InfluxMeasurement).
		AddTag("type", "cowrie.login.success").
		AddField("system", cls.System).
		AddField("eventid", cls.Eventid).
		AddField("username", cls.Username).
		AddField("password", cls.Password).
		AddField("message", cls.Message).
		AddField("time", cls.Time).
		AddField("sensor", cls.Sensor).
		AddField("srcip", cls.SrcIP).
		AddField("session", cls.Session).
		AddField("city", geo.City).
		AddField("country", geo.Country).
		AddField("region", geo.Region).
		AddField("location", geo.Location).
		AddField("latitude", geo.Latitude).
		AddField("longitude", geo.Longitude).
		SetTime(cls.Timestamp)

	write := conn.client.WriteAPIBlocking(getOrganization(), getBucket())
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}

func (conn *InfluxConn) writeCowrieLoginFailed(clf CowrieLoginFailed, geo GeoData) {
	p := influxdb2.NewPointWithMeasurement(settings.InfluxMeasurement).
		AddTag("type", "cowrie.login.failed").
		AddField("system", clf.System).
		AddField("eventid", clf.Eventid).
		AddField("username", clf.Username).
		AddField("password", clf.Password).
		AddField("message", clf.Message).
		AddField("time", clf.Time).
		AddField("sensor", clf.Sensor).
		AddField("srcip", clf.SrcIP).
		AddField("session", clf.Session).
		AddField("city", geo.City).
		AddField("country", geo.Country).
		AddField("region", geo.Region).
		AddField("location", geo.Location).
		AddField("latitude", geo.Latitude).
		AddField("longitude", geo.Longitude).
		SetTime(clf.Timestamp)

	write := conn.client.WriteAPIBlocking(getOrganization(), getBucket())
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}
