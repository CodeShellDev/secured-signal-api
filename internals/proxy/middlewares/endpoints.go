package middlewares

import (
	"net/http"
	"regexp"
	"slices"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
)

var Endpoints Middleware = Middleware{
	Name: "Endpoints",
	Use: endpointsHandler,
}

func endpointsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		endpoints := conf.SETTINGS.ACCESS.ENDPOINTS.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.ENDPOINTS)

		reqPath := req.URL.Path

		if isBlocked(reqPath, matchesPattern, endpoints) {
			logger.Warn("Client tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func matchesPattern(endpoint, pattern string) bool {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return endpoint == pattern
	}

	return re.MatchString(endpoint)
}

func isBlocked(test string, matchFunc func(test, try string) bool, allowBlockSlice structure.AllowBlockSlice) bool {
	if len(allowBlockSlice.Allow) == 0 && len(allowBlockSlice.Block) == 0 {
		// default: allow all
		return false
	}

	isExplicitlyAllowed := slices.ContainsFunc(allowBlockSlice.Allow, func(try string) bool {
		return matchFunc(test, try)
	})
	isExplicitlyBlocked := slices.ContainsFunc(allowBlockSlice.Block, func(try string) bool {
		return matchFunc(test, try)
	})

	// explicit allow > block
	if isExplicitlyAllowed {
		return false
	}
	
	if isExplicitlyBlocked {
		return true
	}

	// allows -> default deny
	if len(allowBlockSlice.Allow) > 0 {
		return true
	}
	
	// only blocks -> default allow
	if len(allowBlockSlice.Block) > 0 {
		return false
	}

	// safety net -> block
	return true
}
