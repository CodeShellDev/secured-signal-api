package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var Endpoints Middleware = Middleware{
	Name: "Endpoints",
	Use: endpointsHandler,
}

func endpointsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		endpoints := conf.SETTINGS.ACCESS.ENDPOINTS.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.ENDPOINTS)

		reqPath := req.URL.Path

		blocked, err := isEndpointBlocked(reqPath, endpoints.Allowed, endpoints.Blocked)

		if err != nil {
			logger.Error("Error during blocked endpoint check: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return 
		}

		if blocked {
			logger.Warn("Client tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func isEndpointBlocked(endpoint string, allowed structure.StringMatchList, blocked structure.StringMatchList) (bool, error) {
	isExplicitlyAllowed, err := allowed.Match(endpoint)

	if err != nil {
		return true, err
	}

	isExplicitlyBlocked, err := blocked.Match(endpoint)

	if err != nil {
		return true, err
	}

	return checkBlockLogic(isExplicitlyAllowed, isExplicitlyBlocked, allowed, blocked), nil
}

func checkBlockLogic[T any](explicitlyAllowed, explicitlyBlocked bool, allowed, blocked []T) bool {
	if len(allowed) == 0 && len(blocked) == 0 {
		// default: allow all
		return false
	}

	// explicit allow > block
	if explicitlyAllowed {
		return false
	}

	if explicitlyBlocked {
		return true
	}

	// allows exist -> default deny
	if len(allowed) > 0 {
		return true
	}

	// only blocks -> default allow
	if len(blocked) > 0 {
		return false
	}

	// safety net -> block
	return true
}