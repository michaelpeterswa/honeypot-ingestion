package logging

import (
	"go.uber.org/zap"
)

func InitLogger(logLevel string) (*zap.Logger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	if logLevel == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}

	defer logger.Sync()
	logger.Info("honeypot-ingestion is initializing...", zap.String("logLevel", logLevel))

	return logger, nil
}
