package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
)

var RequestLogger Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

const loggerKey contextKey = "logger"

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		ip := getContext[net.IP](req, clientIPKey)

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
		conf := getConfigByReq(req)

		logLevel := conf.SERVICE.LOG_LEVEL

		l := logger.Get()

		if strings.TrimSpace(logLevel) != "" {
			l = logger.Get().Sub(logLevel)
		}

		req = setContext(req, loggerKey, l)

		next.ServeHTTP(w, req)
	})
}