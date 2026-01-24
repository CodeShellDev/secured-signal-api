package middlewares

import (
	"net/http"
	"regexp"
	"slices"
	"strings"
)

var Endpoints Middleware = Middleware{
	Name: "Endpoints",
	Use: endpointsHandler,
}

func endpointsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		endpoints := conf.SETTINGS.ACCESS.ENDPOINTS

		if endpoints == nil {
			endpoints = getConfig("").SETTINGS.ACCESS.ENDPOINTS
		}

		reqPath := req.URL.Path

		if isEndpointBlocked(reqPath, endpoints) {
			logger.Warn("Client tried to access blocked endpoint: ", reqPath)
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
	re, err := regexp.Compile(pattern)

	if err != nil {
		return endpoint == pattern
	}

	return re.MatchString(endpoint)
}

func isEndpointBlocked(endpoint string, endpoints []string) bool {
	if len(endpoints) == 0 || endpoints == nil {
		// default: allow all
		return false
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

	// allow rules -> default deny
	if len(allowed) > 0 {
		return true
	}
	
	// only block rules -> default allow
	if len(blocked) > 0 {
		return false
	}

	// safety net -> block
	return true
}
