package middlewares

import (
	"net/http"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type LogMiddleware struct {
	Next http.Handler
}

func (data LogMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Info(req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)

		next.ServeHTTP(w, req)
	})
}
