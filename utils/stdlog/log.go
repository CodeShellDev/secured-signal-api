package stdlog

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/codeshelldev/gotl/pkg/ioutils"
	"github.com/codeshelldev/gotl/pkg/logger"
)

type logLevel int

const logLevelPrefix = "logLevel."

func (l logLevel) String() string {
	return logLevelPrefix + strconv.Itoa(int(l))
}

const (
	FATAL logLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

func normalizeMessage(msg string) string {
	msg = strings.TrimSuffix(msg, "\n")

	msg = strings.ToUpper(msg[:1]) + msg[1:]

	return msg
}

var writer = &ioutils.InterceptWriter{
	Writer: &bytes.Buffer{},
	Hook: func(bytes []byte) {
		msg := string(bytes)
		if len(msg) == 0 {
			return
		}

		level, _ := strconv.Atoi(msg[len(logLevelPrefix):len(logLevelPrefix) + 1])
		msg = msg[len(logLevelPrefix) + 1:]

		msg = normalizeMessage(msg)

		switch (logLevel(level)) {
		case FATAL:
			logger.Fatal(msg)
		case ERROR:
			logger.Error(msg)
		case WARN:
			logger.Warn(msg)
		case INFO:
			logger.Info(msg)
		case DEBUG:
			logger.Debug(msg)
		default:
			logger.Info(msg)
		}
	},
}

var FatalLog *log.Logger = log.New(writer, FATAL.String(), 0)
var ErrorLog *log.Logger = log.New(writer, ERROR.String(), 0)
var WarnLog *log.Logger = log.New(writer, WARN.String(), 0)

var InfoLog *log.Logger = log.New(writer, INFO.String(), 0)
var DebugLog *log.Logger = log.New(writer, DEBUG.String(), 0)
