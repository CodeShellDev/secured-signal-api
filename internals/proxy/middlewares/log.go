package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"go.uber.org/zap/zapcore"
)

var RequestLogger Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

const loggerKey contextKey = "logger"

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := getConfigByReq(req)

		logLevel := conf.SERVICE.LOG_LEVEL

		if strings.TrimSpace(logLevel) == "" {
			logLevel = getConfig("").SERVICE.LOG_LEVEL
		}

		options := logger.DefaultOptions()
		options.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			var name string
			
			if strings.TrimSpace(conf.NAME) != "" {
				name = " " + conf.NAME
			}

			enc.AppendString(caller.TrimmedPath() + name)
		}

		l, err := logger.New(logLevel, options)

		if err != nil {
			logger.Error("Could not create Middleware Logger: ", err.Error())
		}

		if l == nil {
			l = logger.Get()
		}

		req = setContext(req, loggerKey, l)

		ip := getContext[net.IPNet](req, clientIPKey)

		if !l.IsDev() {
			l.Info(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
		} else {
			body, _ := request.GetReqBody(req)

			if body.Data != nil && !body.Empty {
				l.Dev(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery, body.Data)
			} else {
				l.Info(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
			}
		}

		next.ServeHTTP(w, req)
	})
}