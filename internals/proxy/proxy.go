package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/codeshelldev/gotl/pkg/logger"
	m "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	"github.com/codeshelldev/secured-signal-api/utils/stdlog"
)

type Proxy struct {
	Use func() *httputil.ReverseProxy
}

func Create(targetUrl *url.URL) Proxy {
	if targetUrl == nil {
		logger.Fatal("Missing API URL")
		return Proxy{Use: func() *httputil.ReverseProxy {return nil}}
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ModifyResponse = func(res *http.Response) error {
		res.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private, proxy-revalidate")
		res.Header.Set("Pragma", "no-cache")
		res.Header.Set("Expires", "0")
		res.Header.Set("Vary", "*")
		res.Header.Set("Referrer-Policy", "no-referrer")

		return nil
	}

	proxy.ErrorLog = stdlog.ErrorLog

	director := proxy.Director

	proxy.Director = func(req *http.Request) {
		director(req)

		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Host = targetUrl.Host
	}

	return Proxy{Use: func() *httputil.ReverseProxy {return proxy}}
}

func (proxy Proxy) Init() http.Handler {
	handler := m.NewChain().
		Use(m.InternalAPI).
		Use(m.Auth).
		Use(m.InternalMiddlewareLogger).
		Use(m.InternalProxy).
		Use(m.InternalClientIP).
		Use(m.RequestLogger).
		Use(m.InternalAuthRequirement).
		Use(m.Port).
		Use(m.Hostname).
		Use(m.IPFilter).
		Use(m.RateLimit).
		Use(m.Template).
		Use(m.Endpoints).
		Use(m.Mapping).
		Use(m.Policy).
		Use(m.Message).
		Then(proxy.Use())

	return handler
}
