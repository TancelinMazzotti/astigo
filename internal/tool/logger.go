package tool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
)

var (
	Logger *zap.Logger
	once   sync.Once
)

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

func InitLogger(c LoggerConfig) error {
	var err error
	once.Do(func() {
		cfg := zap.Config{
			Encoding:         c.Encoding,
			Level:            zap.NewAtomicLevelAt(parseLevel(c.Level)),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.RFC3339TimeEncoder,
				EncodeDuration: zapcore.MillisDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

		Logger, err = cfg.Build()
	})
	return err
}

func parseLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
