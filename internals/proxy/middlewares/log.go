package middlewares

import (
	"net/http"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
)

var Logging Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !logger.IsDev() {
			logger.Info(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
		} else {
			body, _ := request.GetReqBody(req)

			if body.Data != nil && !body.Empty {
				logger.Dev(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery, body.Data)
			} else {
				logger.Info(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
			}
		}

		next.ServeHTTP(w, req)
	})
}
