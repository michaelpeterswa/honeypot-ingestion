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
		switch results["eventid"] {
		case "cowrie.login.success":
			var cls CowrieLoginSuccess
			json.Unmarshal([]byte(val), &cls)
			fmt.Println(cls.Eventid, cls.Username, cls.Password)
		case "cowrie.login.failed":
			fmt.Println("CowrieLoginFailed")
		case "cowrie.session.connect":
			fmt.Println("CowrieSessionConnect")
		case "cowrie.session.params":
			fmt.Println("CowrieSessionParams")
		case "cowrie.session.closed":
			fmt.Println("CowrieSessionClosed")
		case "cowrie.session.file_download":
			fmt.Println("CowrieSessionFileDownload")
		case "cowrie.command.input":
			fmt.Println("CowrieCommandInput")
		case "cowrie.command.failed":
			fmt.Println("CowrieCommandFailed")
		case "cowrie.direct-tcpip.request":
			fmt.Println("CowrieDirectTCPIPRequest")
		case "cowrie.direct-tcpip.data":
			fmt.Println("CowrieDirectTCPIPData")
		case "cowrie.client.fingerprint":
			fmt.Println("CowrieClientFingerprint")
		case "cowrie.client.kex":
			fmt.Println("CowrieClientKex")
		case "cowrie.client.version":
			fmt.Println("CowrieClientVersion")
		case "cowrie.log.closed":
			fmt.Println("CowrieLogClosed")
		default:
			fmt.Println("EventID not supported...")
		}
	}
}
