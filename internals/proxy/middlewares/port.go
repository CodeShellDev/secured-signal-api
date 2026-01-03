package middlewares

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
)

var Port Middleware = Middleware{
	Name: "Port",
	Use: portHandler,
}

func portHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := getConfigByReq(req)

		allowedPort := conf.SERVICE.PORT

		if strings.TrimSpace(allowedPort) == "" {
			next.ServeHTTP(w, req)
			return
		}

		port := getPort(req.URL)

		if port == "" {
			logger.Error("Could not get Port: invalid scheme")
			http.Error(w, "Bad Request: invalid scheme", http.StatusBadRequest)
			return
		}

		if port != allowedPort {
			logger.Warn("User tried using Token on wrong Port")
			onUnauthorized(w)
			return
		}
	})
}

func getPort(url *url.URL) string {
 	port := url.Port()

	if port == "" {
		return port 
	}

	switch url.Scheme {
	case "https":
		return "443"
	case "http":
		return "80"
	default:
		return ""
	}
}