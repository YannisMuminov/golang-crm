package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func Init(env string) error {
	var cfg zap.Config

	if env == "development" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	var err error

	L, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	if L != nil {
		_ = L.Sync()
	}
}
