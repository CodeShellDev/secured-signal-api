package middlewares

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/codeshelldev/secured-signal-api/internals/config"
)

var Hostname Middleware = Middleware{
	Name: "Hostname",
	Use: hostnameHandler,
}

func hostnameHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		hostnames := conf.SERVICE.HOSTNAMES.OptOrEmpty(config.DEFAULT.SERVICE.HOSTNAMES)

		if len(hostnames) > 0 {
			URL := getContext[*url.URL](req, originURLKey)

			hostname := URL.Hostname()

			if hostname == "" {
				logger.Error("Encountered empty hostname")
				http.Error(w, "Bad Request: invalid hostname", http.StatusBadRequest)
				return
			}

			if !slices.Contains(hostnames, hostname) {
				logger.Warn("Client tried using Token with wrong hostname")
				onUnauthorized(w)
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}