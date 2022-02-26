package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func DefaultProductionLogger(lvl zap.AtomicLevel) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeTime = timeEncoder
	return cfg.Build()
}

func timeEncoder(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(time.UTC().Format("2006/01/02 15:04:05"))
}
