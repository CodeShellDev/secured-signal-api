package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var RequestLogger Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		ip := GetContext[net.IP](req, ClientIPKey)

		if !logger.IsDev() {
			logger.Info(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
		} else {
			body, _ := request.GetReqBody(req)

			if body.Data != nil && !body.Empty {
				logger.Dev(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery, body.Data)
			} else {
				logger.Info(ip.String(), " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
			}
		}

		next.ServeHTTP(w, req)
	})
}

var InternalMiddlewareLogger Middleware = Middleware{
	Name: "_Middleware_Logger",
	Use: middlewareLoggerHandler,
}

func middlewareLoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := GetConfigWithoutDefaultByReq(req)

		var logLevel string

		if conf != nil && conf.TYPE != structure.MAIN {
			logLevel = conf.SERVICE.LOG_LEVEL
		}

		l := logger.Get()

		if strings.TrimSpace(logLevel) != "" {
			l = logger.Get().Sub(logLevel)

			l.SetTransform(func(content string) string {
				return conf.NAME + "\t" + content
			})
		}

		req = SetContext(req, LoggerKey, l)

		next.ServeHTTP(w, req)
	})
}