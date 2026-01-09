package middlewares

import (
	"net"
	"net/http"
)

var InternalClientIP Middleware = Middleware{
	Name: "_Client_IP",
	Use: clientIPHandler,
}

var trustedClientKey contextKey = "isClientTrusted"

func clientIPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		rawTrustedIPs := conf.SETTINGS.ACCESS.TRUSTED_IPS

		if rawTrustedIPs == nil {
			rawTrustedIPs = getConfig("").SETTINGS.ACCESS.TRUSTED_IPS
		}

		ip := getContext[net.IP](req, clientIPKey)

		trustedIPs := parseIPsAndIPNets(rawTrustedIPs)
		trusted := isIPInList(ip, trustedIPs)

		if trusted {
			logger.Dev("Connection from trusted Client: ", ip.String())
		}

		req = setContext(req, trustedClientKey, trusted)

		next.ServeHTTP(w, req)
	})
}