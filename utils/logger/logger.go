package logger

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _logLevel = ""

var logger *zap.Logger

func Init(level string) {
	_logLevel = strings.ToLower(level)

	logLevel := ParseLevel(_logLevel)

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling:    nil,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    CustomEncodeLevel,
			EncodeTime:     zapcore.TimeEncoderOfLayout("02.01 15:04"),
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error

	logger, err = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))

	if err != nil {
		fmt.Println("Encountered Error during Logger Init: ", err.Error())
	}
}

func Format(data ...any) string {
	res := ""

	for _, item := range data {
		switch value := item.(type) {
		case string:
			res += value
		case int:
			res += strconv.Itoa(value)
		case bool:
			if value {
				res += "true"
			} else {
				res += "false"
			}
		default:
			lines := strings.Split(jsonutils.Pretty(value), "\n")

			lineStr := ""

			for _, line := range lines {
				lineStr += "\n" + startColor(color.RGBA{ R: 0, G: 135, B: 95,}) + line + endColor()
			}
			res += lineStr
		}
	}

	return res
}

func Level() string {
	return LevelString(logger.Level())
}

func Sync() {
	logger.Sync()
}
