package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	m "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
)

type Proxy struct {
	Use func() *httputil.ReverseProxy
}

func Create(targetUrl string) Proxy {
	url, _ := url.Parse(targetUrl)

	proxy := httputil.NewSingleHostReverseProxy(url)

	director := proxy.Director

	proxy.Director = func(req *http.Request) {
		director(req)

		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Host = url.Host
	}

	return Proxy{Use: func() *httputil.ReverseProxy {return proxy}}
}

func (proxy Proxy) Init() http.Handler {
	handler := m.NewChain().
		Use(m.Logging).
		Use(m.Server).
		Use(m.Auth).
		Use(m.Endpoints).
		Use(m.Template).
		Use(m.Aliases).
		Use(m.Message).
		Then(proxy.Use())

	return handler
}