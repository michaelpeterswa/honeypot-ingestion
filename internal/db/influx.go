package db

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/geo"
	"github.com/michaelpeterswa/honeypot-ingestion/internal/structs"
	"go.uber.org/zap"
)

type InfluxConn struct {
	client influxdb2.Client
}

func InitInflux(logger *zap.Logger, settings structs.Settings) *InfluxConn {
	logger.Info("Creating InfluxDB v2 connection...")
	client := influxdb2.NewClient(settings.InfluxAddress, settings.InfluxToken)
	return &InfluxConn{client}
}

func (conn *InfluxConn) WriteCowrieLoginSuccess(logger *zap.Logger, settings structs.Settings, cls structs.CowrieLoginSuccess, geo geo.GeoData) {
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
		SetTime(time.Now())

	write := conn.client.WriteAPIBlocking(settings.InfluxOrganization, settings.InfluxBucket)
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}

func (conn *InfluxConn) WriteCowrieLoginFailed(logger *zap.Logger, settings structs.Settings, clf structs.CowrieLoginFailed, geo geo.GeoData) {
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
		SetTime(time.Now())

	write := conn.client.WriteAPIBlocking(settings.InfluxOrganization, settings.InfluxBucket)
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}

func (conn *InfluxConn) WriteCowrieSessionConnect(logger *zap.Logger, settings structs.Settings, csc structs.CowrieSessionConnect, geo geo.GeoData) {
	p := influxdb2.NewPointWithMeasurement(settings.InfluxMeasurement).
		AddTag("type", "cowrie.session.connect").
		AddField("system", csc.System).
		AddField("eventid", csc.Eventid).
		AddField("srcip", csc.SrcIP).
		AddField("srcport", csc.SrcPort).
		AddField("dstip", csc.DstIP).
		AddField("dstport", csc.DstPort).
		AddField("session", csc.Session).
		AddField("protocol", csc.Protocol).
		AddField("message", csc.Message).
		AddField("time", csc.Time).
		AddField("sensor", csc.Sensor).
		AddField("city", geo.City).
		AddField("country", geo.Country).
		AddField("region", geo.Region).
		AddField("location", geo.Location).
		AddField("latitude", geo.Latitude).
		AddField("longitude", geo.Longitude).
		SetTime(time.Now())

	write := conn.client.WriteAPIBlocking(settings.InfluxOrganization, settings.InfluxBucket)
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}

func (conn *InfluxConn) WriteCowrieCommandInput(logger *zap.Logger, settings structs.Settings, cci structs.CowrieCommandInput, geo geo.GeoData) {
	p := influxdb2.NewPointWithMeasurement(settings.InfluxMeasurement).
		AddTag("type", "cowrie.command.input").
		AddField("system", cci.System).
		AddField("eventid", cci.Eventid).
		AddField("input", cci.Input).
		AddField("message", cci.Message).
		AddField("time", cci.Time).
		AddField("sensor", cci.Sensor).
		AddField("timestamp", cci.Timestamp).
		AddField("srcip", cci.SrcIP).
		AddField("session", cci.Session).
		AddField("city", geo.City).
		AddField("country", geo.Country).
		AddField("region", geo.Region).
		AddField("location", geo.Location).
		AddField("latitude", geo.Latitude).
		AddField("longitude", geo.Longitude).
		SetTime(time.Now())

	write := conn.client.WriteAPIBlocking(settings.InfluxOrganization, settings.InfluxBucket)
	err := write.WritePoint(context.Background(), p)
	if err != nil {
		logger.Error("Write to InfluxDB failed...", zap.Error(err))
	}
}
