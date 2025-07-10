package middlewares

import (
	"net/http"
	"slices"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type EndpointsMiddleware struct {
	Next             http.Handler
	BlockedEndpoints []string
}

func (data EndpointsMiddleware) Use() http.Handler {
	next := data.Next
	BLOCKED_ENDPOINTS := data.BlockedEndpoints

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqPath := req.URL.Path

		if slices.Contains(BLOCKED_ENDPOINTS, reqPath) {
			log.Warn("User tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}
