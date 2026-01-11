package middlewares

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var InternalProxy Middleware = Middleware{
	Name: "_Proxy",
	Use: proxyHandler,
}

const trustedProxyKey contextKey = "isProxyTrusted"
const clientIPKey contextKey = "clientIP"
const originURLKey contextKey = "originURL"

func proxyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)
		
		conf := getConfigByReq(req)

		rawTrustedProxies := conf.SETTINGS.ACCESS.TRUSTED_PROXIES

		if rawTrustedProxies == nil {
			rawTrustedProxies = getConfig("").SETTINGS.ACCESS.TRUSTED_PROXIES
		}

		var trusted bool
		var ip net.IP

		host, _, _ := net.SplitHostPort(req.RemoteAddr)

		originUrl := req.Proto + "://" + req.URL.Host

		ip = net.ParseIP(host)

		if len(rawTrustedProxies) != 0 {
			trustedProxies := parseIPsAndIPNets(rawTrustedProxies)

			trusted = isIPInList(ip, trustedProxies)
		}

		if trusted {
			realIP, err := getRealIP(req)

			if err != nil {
				logger.Error("Could not get real IP: ", err.Error())
			}

			if realIP != nil {
				ip = realIP
			}

			XFHost := req.Header.Get("X-Forwarded-Host")
			XFProto := req.Header.Get("X-Forwarded-Proto")
			XFPort := req.Header.Get("X-Forwarded-Port")

			if XFHost == "" || XFProto == "" || XFPort == "" {
				logger.Warn("Missing X-Forwarded-* headers")
			}

			originUrl = XFProto + "://" + XFHost + ":" + XFPort
		}

		originURL, err := url.Parse(originUrl)

		if err != nil {
			logger.Error("Could not parse Url: ", originUrl)
			http.Error(w, "Bad Request: invalid Url", http.StatusBadRequest)
			return
		}

		req = setContext(req, trustedProxyKey, trusted)
		req = setContext(req, originURLKey, originURL)

		req = setContext(req, clientIPKey, ip)

		next.ServeHTTP(w, req)
	})
}

func parseIP(str string) (*net.IPNet, error) {
    if !strings.Contains(str, "/") {
        ip := net.ParseIP(str)

        if ip == nil {
            return nil, errors.New("invalid ip: " + str)
        }

        var mask net.IPMask

        if ip.To4() != nil {
            mask = net.CIDRMask(32, 32) // IPv4 /32
        } else {
            mask = net.CIDRMask(128, 128) // IPv6 /128
        }
		
        return &net.IPNet{IP: ip, Mask: mask}, nil
    }

    ip, network, err := net.ParseCIDR(str)
    if err != nil {
        return nil, err
    }

    if !ip.Equal(network.IP) {
        var mask net.IPMask

        if ip.To4() != nil {
            mask = net.CIDRMask(32, 32) // IPv4 /32
        } else {
            mask = net.CIDRMask(128, 128) // IPv6 /128
        }

        return &net.IPNet{IP: ip, Mask: mask}, nil
    }

    return network, nil
}

func parseIPsAndIPNets(array []string) []*net.IPNet {
	ipNets := []*net.IPNet{}

	for _, item := range array {
		ipNet, err := parseIP(item)

		if err != nil {
			continue
		}

		ipNets = append(ipNets, ipNet)
	}

	return ipNets
}

func getRealIP(req *http.Request) (net.IP, error) {
	XFF := req.Header.Get("X-Forwarded-For")

	if XFF != "" {
		ips := strings.Split(XFF, ",")
		
		realIP := net.ParseIP(strings.TrimSpace(ips[0]))

		if realIP == nil {
			return nil, errors.New("malformed X-Forwarded-For header")
		}

		return realIP, nil
	}

	return nil, errors.New("no X-Forwarded-For header present")
}

func isIPInList(ip net.IP, list []*net.IPNet) bool {
	for _, net := range list {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}