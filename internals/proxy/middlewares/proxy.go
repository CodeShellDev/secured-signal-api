package middlewares

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/utils/netutils"
)

var InternalProxy Middleware = Middleware{
	Name: "_Proxy",
	Use: proxyHandler,
}

const trustedProxyKey contextKey = "isProxyTrusted"
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

		rawTrustedProxies := conf.SETTINGS.ACCESS.TRUSTED_PROXIES.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.TRUSTED_PROXIES)

        var trusted bool
		var ip net.IP

		host, _, _ := net.SplitHostPort(req.RemoteAddr)

        originUrl := parseOrigin(req.Proto, req.Host)

		ip = net.ParseIP(host)

		if len(rawTrustedProxies) != 0 {
            trustedProxies := parseIPsAndNets(rawTrustedProxies)
			trusted = netutils.IsIPIn(ip, trustedProxies)
		}

		if trusted {
            var forwardedEntries []ForwardedEntry

            if req.Header.Get("Forwarded") != "" {
				forwardedEntries = parseForwarded(req.Header.Get("Forwarded"))
			} else {
				forwardedEntries = parseXForwardedHeaders(req.Header)
			}

            if len(forwardedEntries) != 0 {
                ip = parseForIP(forwardedEntries[0].For)

                originUrl = parseOrigin(forwardedEntries[0].Proto, forwardedEntries[0].Host)
            }
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

func parseOrigin(proto, host string) string {
    if !strings.Contains(host, ":") {
        if proto == "https" {
            host += ":443"
        } else {
            host += ":80"
        }
    }

    return proto + "://" + host
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