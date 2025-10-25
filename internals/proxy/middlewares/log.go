package middlewares

import (
	"net/http"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var Logging Middleware = Middleware{
	Name: "Logging",
	Use: loggingHandler,
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Info(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)

		next.ServeHTTP(w, req)
	})
}
