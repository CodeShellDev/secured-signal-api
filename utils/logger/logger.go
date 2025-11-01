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

var _log *zap.Logger
var _logLevel = ""

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

	_log, err = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))

	if err != nil {
		fmt.Println("Encountered Error during Log.Init(): ", err.Error())
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
		default:
			res += "\n" + ColorCode(jsonutils.Pretty(value), color.RGBA{
				R: 0, G: 215, B: 135,
			})
		}
	}

	return res
}

func Level() string {
	return LevelString(_log.Level())
}

func Info(data ...any) {
	_log.Info(Format(data...))
}

func Debug(data ...any) {
	_log.Debug(Format(data...))
}

func Dev(data ...any) {
	ok := _log.Check(DeveloperLevel, Format(data...))

	if ok != nil {
		ok.Write()
	}
}

func Error(data ...any) {
	_log.Error(Format(data...))
}

func Fatal(data ...any) {
	_log.Fatal(Format(data...))
}

func Warn(data ...any) {
	_log.Warn(Format(data...))
}

func Sync() {
	_ = _log.Sync()
}
