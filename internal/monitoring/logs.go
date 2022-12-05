package monitoring

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func NewLogger(loglevel string) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	err := cfg.Level.UnmarshalText([]byte(loglevel))
	if err != nil {
		log.Panicf("failed at configuring logger: %v", err)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("failed at configuring logger: %v", err)
	}
	return logger.Sugar()
}
