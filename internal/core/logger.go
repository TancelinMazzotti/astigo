package core

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

func NewLogger(c LoggerConfig) (*zap.Logger, error) {
	var err error

	cfg := zap.Config{
		Encoding:         c.Encoding,
		Level:            zap.NewAtomicLevelAt(parseLevel(c.Level)),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "msg",
			LevelKey:         "level",
			TimeKey:          "time",
			NameKey:          "logger",
			CallerKey:        "caller",
			StacktraceKey:    "stacktrace",
			SkipLineEnding:   false,
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.RFC3339TimeEncoder,
			EncodeDuration:   zapcore.MillisDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: " ",
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
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
