package middlewares

import (
	"net"
	"net/http"
	"slices"
	"strings"
)

var IPFilter Middleware = Middleware{
	Name: "IP Filter",
	Use: ipFilterHandler,
}

func ipFilterHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		ipFilter := conf.SETTINGS.ACCESS.IP_FILTER

		if ipFilter == nil {
			ipFilter = getConfig("").SETTINGS.ACCESS.ENDPOINTS
		}

		ip := getContext[net.IP](req, clientIPKey)

		if isIPBlocked(ip, ipFilter) {
			logger.Warn("Client IP is blocked by filter: ", ip.String())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

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

func isIPBlocked(ip net.IP, ipfilter []string) (bool) {
	if len(ipfilter) == 0 || ipfilter == nil {
		// default: allow all
		return false
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
