package middlewares

import (
	"net/http"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/codeshelldev/secured-signal-api/utils/request"
)

var Logging Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !log.IsDev() {
			log.Info(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
		} else {
			body, _ := request.GetReqBody(req)

			log.Dev(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery, body)
		}

		next.ServeHTTP(w, req)
	})
}
