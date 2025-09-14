package middlewares

import (
	"net/http"
	"slices"
	"strings"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type EndpointsMiddleware struct {
	Next             http.Handler
}

func (data EndpointsMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		blockedEndpoints := settings.BLOCKED_ENDPOINTS
		allowedEndpoints := settings.ALLOWED_ENDPOINTS

		if blockedEndpoints == nil {
			blockedEndpoints = getSettings("*").BLOCKED_ENDPOINTS
		}

		reqPath := req.URL.Path

		if isBlocked(reqPath, allowedEndpoints, blockedEndpoints) {
			log.Warn("User tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func isBlocked(endpoint string, allowed []string, blocked []string) bool {
	var result bool

	if blocked == nil {
		return false
	}

	if allowed == nil {
		return true
	}

	isBlocked := slices.ContainsFunc(blocked, func(try string) bool {
		return strings.HasPrefix(endpoint, try)
	})

	isExplictlyAllowed := slices.ContainsFunc(allowed, func(try string) bool {
		return strings.HasPrefix(endpoint, try)
	})

	result = isBlocked && !isExplictlyAllowed

	return result
}