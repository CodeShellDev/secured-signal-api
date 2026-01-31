package middlewares

import (
	"net/http"
	"path"
	"slices"
	"strings"

	log "github.com/codeshelldev/gotl/pkg/logger"
)

var Endpoints Middleware = Middleware{
	Name: "Endpoints",
	Use: endpointsHandler,
}

func endpointsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := getConfigByReq(req)

		endpoints := conf.SETTINGS.ACCESS.ENDPOINTS

		if endpoints == nil {
			endpoints = getConfig("").SETTINGS.ACCESS.ENDPOINTS
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

func matchesPattern(endpoint, pattern string) bool {
	ok, _ := path.Match(pattern, endpoint)
	return ok
}

func isBlocked(endpoint string, endpoints []string) bool {
	if len(endpoints) == 0 {
		// default: block all
		return true
	}

	allowed, blocked := getEndpoints(endpoints)

	isExplicitlyAllowed := slices.ContainsFunc(allowed, func(try string) bool {
		return matchesPattern(endpoint, try)
	})
	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try string) bool {
		return matchesPattern(endpoint, try)
	})

	// explicit allow > block
	if isExplicitlyAllowed {
		return false
	}
	
	if isExplicitlyBlocked {
		return true
	}

	// only allowed endpoints -> block anything not allowed
	if len(allowed) > 0 && len(blocked) == 0 {
		return true
	}

	// only blocked endpoints -> allow anything not blocked
	if len(blocked) > 0 && len(allowed) == 0 {
		return false
	}

	// no match -> default: block all
	return true
}
