package middlewares

import (
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
)

var IPFilter Middleware = Middleware{
	Name: "IP Filter",
	Use: ipFilterHandler,
}

var trustedClientKey contextKey = "isClientTrusted"

func ipFilterHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := getConfigByReq(req)

		ipFilter := conf.SETTINGS.ACCESS.IP_FILTER

		logger.Dev(ipFilter)

		if ipFilter == nil {
			ipFilter = getConfig("").SETTINGS.ACCESS.ENDPOINTS
		}

		ip := getContext[net.IP](req, clientIPKey)

		block, trusted := blockIPOrTrust(ip, ipFilter)

		logger.Dev(block, trusted)

		if block {
			logger.Warn("Client IP is blocked by filter: ", ip.String())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if trusted {
			logger.Dev("Connection from trusted Client: ", ip.String())
		}

		req = setContext(req, trustedClientKey, trusted)

		next.ServeHTTP(w, req)
	})
}

func getIPNets(ipNets []string) ([]string, []string) {
	blockedIPNets := []string{}
	allowedIPNets := []string{}

	for _, ipNet := range ipNets {
		ip, block := strings.CutPrefix(ipNet, "!")

		if block {
			blockedIPNets = append(blockedIPNets, ip)
		} else {
			allowedIPNets = append(allowedIPNets, ip)
		}
	}

	return allowedIPNets, blockedIPNets
}

func blockIPOrTrust(ip net.IP, ipfilter []string) (bool, bool) {
	if len(ipfilter) == 0 {
		// default: allow all, but do not trust
		return false, false
	}

	rawAllowed, rawBlocked := getIPNets(ipfilter)

	allowed := parseIPsAndIPNets(rawAllowed)
	blocked := parseIPsAndIPNets(rawBlocked)

	isExplicitlyAllowed := slices.ContainsFunc(allowed, func(try *net.IPNet) bool {
		return try.Contains(ip)
	})
	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try *net.IPNet) bool {
		return try.Contains(ip)
	})

	// explicit allow > block
	if isExplicitlyAllowed {
		return false, true
	}
	
	if isExplicitlyBlocked {
		return true, false
	}

	// default: allow all, but do not trust
	return false, false
}
