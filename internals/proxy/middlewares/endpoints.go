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

		if blockedEndpoints == nil && allowedEndpoints == nil {
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
	if blocked == nil {
		blocked = []string{}
	}

	if allowed == nil {
		allowed = []string{}
	}

	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try string) bool {
		log.Dev("Checking " + try +  " against " + endpoint)
		return strings.HasPrefix(endpoint, try)
	})

	isExplictlyAllowed := slices.ContainsFunc(allowed, func(try string) bool {
		log.Dev("Checking " + try + " against " + endpoint)
		return strings.HasPrefix(endpoint, try)
	})

	// Block all except explicitly Allowed
	if len(blocked) == 0 && len(allowed) != 0 {
		return !isExplictlyAllowed
	}

	// Allow all except explicitly Blocked
	if len(allowed) == 0 && len(blocked) != 0{
		return isExplicitlyBlocked
	}

	// Excplicitly Blocked except excplictly Allowed
	if len(blocked) != 0 && len(allowed) != 0 {
		return isExplicitlyBlocked && !isExplictlyAllowed
	}

	// Block all
	return true
}