package middlewares

import (
	"net"
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/netutils"
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

		rawTrustedIPs := conf.SETTINGS.ACCESS.TRUSTED_IPS.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.TRUSTED_IPS)

		ip := getContext[net.IP](req, clientIPKey)

		trustedIPs := parseIPsAndNets(rawTrustedIPs)
		trusted := netutils.IsIPIn(ip, trustedIPs)

		if trusted {
			logger.Dev("Connection from trusted Client: ", ip.String())
		}

		req = setContext(req, trustedClientKey, trusted)

		next.ServeHTTP(w, req)
	})
}

func parseIPsAndNets(ipNets []structure.IPOrNet) []*net.IPNet {
    out := []*net.IPNet{}

    for _, ipNet := range ipNets {
        out = append(out, ipNet.IPNet)
    }

    return out
}