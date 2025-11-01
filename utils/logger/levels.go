package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const DeveloperLevel zapcore.Level = -2

func ParseLevel(s string) zapcore.Level {
	switch strings.ToLower(s) {
	case "dev":
		return DeveloperLevel
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func LevelString(l zapcore.Level) string {
	switch l {
	case DeveloperLevel:
		return "dev"
	default:
		return l.CapitalString()
	}
}

func CapitalLevel(l zapcore.Level) string {
	switch l {
	case DeveloperLevel:
		return "DEV"
	default:
		return l.CapitalString()
	}
}

func CustomEncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case DeveloperLevel:
		enc.AppendString(CapitalLevel(l))
	default:
		zapcore.CapitalColorLevelEncoder(l, enc)
	}
}
