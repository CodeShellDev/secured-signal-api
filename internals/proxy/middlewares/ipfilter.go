package middlewares

import (
	"net"
	"net/http"
	"slices"

	g "github.com/codeshelldev/gotl/pkg/configutils/types/generics"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var IPFilter Middleware = Middleware{
	Name: "IP Filter",
	Use: ipFilterHandler,
}

func ipFilterHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		ipFilter := conf.SETTINGS.ACCESS.IP_FILTER.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.IP_FILTER)

		if len(ipFilter.Allowed) == 0 && len(ipFilter.Blocked) == 0 {
			next.ServeHTTP(w, req)
			return
		}

		ip := GetContext[net.IP](req, ClientIPKey)

		if isIPBlocked(ip, ipFilter.Allowed, ipFilter.Blocked) {
			logger.Warn("Client IP is blocked by filter: ", ip.String())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func isIPBlocked(ip net.IP, allowed []g.IPOrNet, blocked []g.IPOrNet) bool {
	isExplicitlyAllowed := slices.ContainsFunc(allowed, func(try g.IPOrNet) bool {
		tryIP := net.IPNet(try)
		
		return tryIP.Contains(ip)
	})

	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try g.IPOrNet) bool {
		tryIP := net.IPNet(try)
		
		return tryIP.Contains(ip)
	})

	return checkBlockLogic(isExplicitlyAllowed, isExplicitlyBlocked, allowed, blocked)
}