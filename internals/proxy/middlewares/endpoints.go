package middlewares

import (
	"net/http"
	"slices"
	"strings"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var Endpoints Middleware = Middleware{
	Name: "Endpoints",
	Use: endpointsHandler,
}

func endpointsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		endpoints := settings.ACCESS.ENDPOINTS

		if endpoints == nil {
			endpoints = getSettings("*").ACCESS.ENDPOINTS
		}

		reqPath := req.URL.Path

		if isBlocked(reqPath, endpoints) {
			log.Warn("User tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getEndpoints(endpoints []string) ([]string, []string) {
	blockedEndpoints := []string{}
	allowedEndpoints := []string{}

	for _, endpoint := range endpoints {
		endpoint, block := strings.CutPrefix(endpoint, "!")

		if block {
			blockedEndpoints = append(blockedEndpoints, endpoint)
		} else {
			allowedEndpoints = append(allowedEndpoints, endpoint)
		}
	}

	return allowedEndpoints, blockedEndpoints
}

func isBlocked(endpoint string, endpoints []string) bool {
	if endpoints == nil {
		return false
	} else if len(endpoints) <= 0 {
		return false
	}

	allowed, blocked := getEndpoints(endpoints)

	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try string) bool {
		return strings.HasPrefix(endpoint, try)
	})

	isExplictlyAllowed := slices.ContainsFunc(allowed, func(try string) bool {
		return strings.HasPrefix(endpoint, try)
	})

	// Block all except explicitly Allowed
	if len(blocked) == 0 && len(allowed) != 0 {
		return !isExplictlyAllowed
	}

	// Allow all except explicitly Blocked
	if len(allowed) == 0 && len(blocked) != 0 {
		return isExplicitlyBlocked
	}

	// Excplicitly Blocked except excplictly Allowed
	if len(blocked) != 0 && len(allowed) != 0 {
		return isExplicitlyBlocked && !isExplictlyAllowed
	}

	// Block all
	return true
}
