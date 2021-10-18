package main

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger(logLevel string) {
	var err error
	if logLevel == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Println("Logger initialization failed...")
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Println("Logger initialization failed...")
		}
	}

	defer logger.Sync()
	logger.Info("honeypot-ingestion is initializing...", zap.String("logLevel", settings.ZapLevel))
}
