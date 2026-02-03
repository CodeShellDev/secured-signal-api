package middlewares

import (
	"net"
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/utils/netutils"
)

var IPFilter Middleware = Middleware{
	Name: "IP Filter",
	Use: ipFilterHandler,
}

func ipFilterHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		ipFilter := conf.SETTINGS.ACCESS.IP_FILTER.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.IP_FILTER)

		logger.Dev(conf.SETTINGS.ACCESS.IP_FILTER)

		ip := getContext[net.IP](req, clientIPKey)

		if isBlocked("", func(_, try string) bool {
			tryIP, err := netutils.ParseIPorNet(try)
			
			return tryIP.Contains(ip) && err == nil
		}, ipFilter) {
			logger.Warn("Client IP is blocked by filter: ", ip.String())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}