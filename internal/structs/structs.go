package structs

import "time"

type Settings struct {
	RedisURL      string `yaml:"redis-url"`
	RedisPort     int    `yaml:"redis-port"`
	RedisPassword string `yaml:"redis-password"`

	CowrieKey string `yaml:"cowrie-key"`

	InfluxAddress      string `yaml:"influx-address"`
	InfluxToken        string `yaml:"influx-token"`
	InfluxOrganization string `yaml:"influx-organization"`
	InfluxBucket       string `yaml:"influx-bucket"`
	InfluxMeasurement  string `yaml:"influx-measurement"`

	IPInfoKey string `yaml:"ipinfo-key"`

	ZapLevel string `yaml:"zap-level"`
}

// cowrie.login.success
// cowrie.login.failed
// cowrie.session.connect
// cowrie.session.params
// cowrie.session.closed
// cowrie.session.file_download
// cowrie.command.input
// cowrie.command.failed
// cowrie.direct-tcpip.request
// cowrie.direct-tcpip.data
// cowrie.client.fingerprint
// cowrie.client.kex
// cowrie.client.version
// cowrie.log.closed

type CowrieClientFingerprint struct {
	System      string    `json:"system"`
	Eventid     string    `json:"eventid"`
	Username    string    `json:"username"`
	Fingerprint string    `json:"fingerprint"`
	Key         string    `json:"key"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Time        float64   `json:"time"`
	Sensor      string    `json:"sensor"`
	Timestamp   time.Time `json:"timestamp"`
	SrcIP       string    `json:"src_ip"`
	Session     string    `json:"session"`
}

type CowrieLogClosed struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Ttylog    string    `json:"ttylog"`
	Size      int       `json:"size"`
	Shasum    string    `json:"shasum"`
	Duplicate bool      `json:"duplicate"`
	Duration  float64   `json:"duration"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieLoginSuccess struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieDirectTCPIPData struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	DstIP     string    `json:"dst_ip"`
	DstPort   int       `json:"dst_port"`
	Data      string    `json:"data"`
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieClientVersion struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Version   string    `json:"version"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieLoginFailed struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieCommandInput struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Input     string    `json:"input"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieClientKex struct {
	System          string    `json:"system"`
	Eventid         string    `json:"eventid"`
	Hassh           string    `json:"hassh"`
	HasshAlgorithms string    `json:"hasshAlgorithms"`
	KexAlgs         []string  `json:"kexAlgs"`
	KeyAlgs         []string  `json:"keyAlgs"`
	EncCS           []string  `json:"encCS"`
	MacCS           []string  `json:"macCS"`
	CompCS          []string  `json:"compCS"`
	LangCS          []string  `json:"langCS"`
	Message         string    `json:"message"`
	Time            float64   `json:"time"`
	Sensor          string    `json:"sensor"`
	Timestamp       time.Time `json:"timestamp"`
	SrcIP           string    `json:"src_ip"`
	Session         string    `json:"session"`
}

type CowrieSessionConnect struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	SrcIP     string    `json:"src_ip"`
	SrcPort   int       `json:"src_port"`
	DstIP     string    `json:"dst_ip"`
	DstPort   int       `json:"dst_port"`
	Session   string    `json:"session"`
	Protocol  string    `json:"protocol"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
}

type CowrieSessionClosed struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Duration  float64   `json:"duration"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieCommandFailed struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Input     string    `json:"input"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieSessionFileDownload struct {
	Eventid   string    `json:"eventid"`
	URL       string    `json:"url"`
	Outfile   string    `json:"outfile"`
	Shasum    string    `json:"shasum"`
	Sensor    string    `json:"sensor"`
	Time      float64   `json:"time"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieSessionParams struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	Arch      string    `json:"arch"`
	Message   []string  `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	Session   string    `json:"session"`
}

type CowrieDirectTCPIPRequest struct {
	System    string    `json:"system"`
	Eventid   string    `json:"eventid"`
	DstIP     string    `json:"dst_ip"`
	DstPort   int       `json:"dst_port"`
	SrcIP     string    `json:"src_ip"`
	SrcPort   int       `json:"src_port"`
	Message   string    `json:"message"`
	Time      float64   `json:"time"`
	Sensor    string    `json:"sensor"`
	Timestamp time.Time `json:"timestamp"`
	Session   string    `json:"session"`
}
