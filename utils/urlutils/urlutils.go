package urlutils

import "net/url"

func NormalizeURL(url *url.URL) string {
    host := url.Hostname()
    port := url.Port()

    if (url.Scheme == "https" && port == "443") ||
       (url.Scheme == "http" && port == "80") ||
       port == "" {
        return url.Scheme + "://" + host
    }

    return url.Scheme + "://" + host + ":" + port
}