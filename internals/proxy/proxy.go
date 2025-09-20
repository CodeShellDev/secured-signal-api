package proxy

import (
	"net/http/httputil"
	"net/url"
)

func Create(targetUrl string) *httputil.ReverseProxy {
	url, _ := url.Parse(targetUrl)

	proxy := httputil.NewSingleHostReverseProxy(url)

	return proxy
}