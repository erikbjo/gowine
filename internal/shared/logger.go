package shared

import (
	"go.uber.org/zap"
)

func CreateSugaredLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
