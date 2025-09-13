package middlewares

import (
	"net/http"
	"slices"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type EndpointsMiddleware struct {
	Next             http.Handler
}

func (data EndpointsMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		blockedEndpoints := GetSettings(req).BLOCKED_ENDPOINTS

		reqPath := req.URL.Path

		if slices.Contains(blockedEndpoints, reqPath) {
			log.Warn("User tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}
