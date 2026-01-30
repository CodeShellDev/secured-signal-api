package middlewares

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

var InternalProxy Middleware = Middleware{
	Name: "_Proxy",
	Use: proxyHandler,
}

const clientIPKey contextKey = "clientIP"
const originURLKey contextKey = "originURL"

type ForwardedEntry struct {
	For		string
	Host	string
	Proto	string
}

type OriginInfo	struct {
	IP		net.IP
	Host	string
	Proto	string
}

func proxyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)
		
		conf := getConfigByReq(req)

		rawTrustedProxies := conf.SETTINGS.ACCESS.TRUSTED_PROXIES

		if rawTrustedProxies == nil {
			rawTrustedProxies = getConfig("").SETTINGS.ACCESS.TRUSTED_PROXIES
		}

		var ip net.IP
		var originUrl string

		if len(rawTrustedProxies) != 0 {
			trustedProxies := parseIPsAndIPNets(rawTrustedProxies)

			var forwardedEntries []ForwardedEntry

			if req.Header.Get("Forwarded") != "" {
				forwardedEntries = parseForwarded(req.Header.Get("Forwarded"))
			} else {
				forwardedEntries = parseXForwardedHeaders(req.Header)
			}

			logger.Dev(forwardedEntries)

			if len(forwardedEntries) != 0 {
				originInfo := getOriginFromForwarded(forwardedEntries, trustedProxies)
				ip = originInfo.IP

                logger.Dev(originInfo)

				originUrl = originInfo.Proto + "://" + originInfo.Host
			}
		}

		if ip == nil {
			host, _, _ := net.SplitHostPort(req.RemoteAddr)

			ip = net.ParseIP(host)
		}

		if originUrl == "" {
			originUrl = req.Proto + "://" + req.Host

			if !strings.Contains(req.Host, ":") {
				if req.Proto == "https" {
					originUrl += ":443"
				} else {
					originUrl += ":80"
				}
			}

            logger.Dev(originUrl)
		}

		originURL, err := url.Parse(originUrl)

		if err != nil {
			logger.Error("Could not parse Url: ", originUrl)
			http.Error(w, "Bad Request: invalid Url", http.StatusBadRequest)
			return
		}

		req = setContext(req, originURLKey, originURL)

		req = setContext(req, clientIPKey, ip)

		next.ServeHTTP(w, req)
	})
}

func getOriginFromForwarded(entries []ForwardedEntry, trusted []*net.IPNet) OriginInfo {
    var origin OriginInfo

	// reverse to place origin client last
	slices.Reverse(entries)

    for _, entry := range entries {
        ip := parseForIP(entry.For)

        if ip == nil {
            continue
        }

		// ip not trusted => use as client ip
        if !isIPInList(ip, trusted) {
            origin.IP = ip
            origin.Proto = entry.Proto
            origin.Host = entry.Host
            break
        }
    }

    return origin
}

func parseForIP(value string) net.IP {
    value = strings.TrimSpace(value)
    value = strings.Trim(value, `"`)
    value = strings.Trim(value, "[]")

	host, _, err := net.SplitHostPort(value)
    if err == nil {
        value = host
    }

    return net.ParseIP(value)
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

    _, network, err := net.ParseCIDR(str)

    if err != nil {
        return nil, err
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

func parseXForwardedHeaders(headers http.Header) []ForwardedEntry {
	var entries []ForwardedEntry

    XFF := headers.Get("X-Forwarded-For")
    if XFF == "" {
        return nil
    }

    parts := strings.Split(XFF, ",")

    XFProto := headers.Get("X-Forwarded-Proto")
    XFHost  := headers.Get("X-Forwarded-Host")

    for i, part := range parts {
        ip := strings.TrimSpace(part)
        if ip == "" {
            continue
        }

        entry := ForwardedEntry{
            For: ip,
        }

        if i == 0 {
            if XFProto != "" {
                entry.Proto = XFProto
            }
            if XFHost != "" {
                entry.Host = XFHost
            }
        }

        entries = append(entries, entry)
    }

    return entries
}

func parseForwarded(header string) []ForwardedEntry {
    var entries []ForwardedEntry

    for part := range strings.SplitSeq(header, ",") {
        entry := ForwardedEntry{}
        params := strings.SplitSeq(part, ";")

        for param := range params {
            keyValuePair := strings.SplitN(strings.TrimSpace(param), "=", 2)

            if len(keyValuePair) != 2 {
                continue
            }

            key := strings.ToLower(keyValuePair[0])
            value := strings.Trim(keyValuePair[1], `"`)

            switch key {
            case "for":
                entry.For = value
            case "proto":
                entry.Proto = value
            case "host":
                entry.Host = value
            }
        }

        entries = append(entries, entry)
    }

    return entries
}

func isIPInList(ip net.IP, list []*net.IPNet) bool {
	for _, net := range list {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}