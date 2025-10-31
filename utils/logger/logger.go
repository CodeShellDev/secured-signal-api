package logger

import (
	"fmt"
	"strings"

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
		fmt.Println("Encountered Error during Log.Init(): err.Error()")
	}
}

func Level() string {
	return _log.Level().String()
}

func Info(msg ...string) {
	_log.Info(strings.Join(msg, ""))
}

func Debug(msg ...string) {
	_log.Debug(strings.Join(msg, ""))
}

func Dev(msg ...string) {
	ok := _log.Check(DeveloperLevel, strings.Join(msg, ""))

	if ok != nil {
		ok.Write()
	}
}

func Error(msg ...string) {
	_log.Error(strings.Join(msg, ""))
}

func Fatal(msg ...string) {
	_log.Fatal(strings.Join(msg, ""))
}

func Warn(msg ...string) {
	_log.Warn(strings.Join(msg, ""))
}

func Sync() {
	_ = _log.Sync()
}
