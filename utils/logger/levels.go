package logger

import "go.uber.org/zap/zapcore"

func Dev(data ...any) {
	ok := logger.Check(DeveloperLevel, Format(data...))

	if ok != nil {
		ok.Write()
	}
}

func Debug(data ...any) {
	logger.Debug(Format(data...))
}

func Info(data ...any) {
	logger.Info(Format(data...))
}

func Warn(data ...any) {
	logger.Warn(Format(data...))
}

func Error(data ...any) {
	logger.Error(Format(data...))
}

func Fatal(data ...any) {
	logger.Fatal(Format(data...))
}

func IsDev() bool {
	return logger.Level().Enabled(DeveloperLevel)
}

func IsDebug() bool {
	return logger.Level().Enabled(zapcore.DebugLevel)
}
func IsInfo() bool {
	return logger.Level().Enabled(zapcore.InfoLevel)
}
func IsWarn() bool {
	return logger.Level().Enabled(zapcore.WarnLevel)
}
func IsError() bool {
	return logger.Level().Enabled(zapcore.ErrorLevel)
}
func IsFatal() bool {
	return logger.Level().Enabled(zapcore.FatalLevel)
}