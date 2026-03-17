package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/codeshelldev/gotl/pkg/logger"
	m "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
)

type Proxy struct {
	Use func() *httputil.ReverseProxy
}

func Create(targetUrl *url.URL) Proxy {
	if targetUrl == nil {
		logger.Fatal("Missing API URL")
		return Proxy{Use: func() *httputil.ReverseProxy {return nil}}
	}

	modifyResponse := m.NewResponseChain().
		Use(m.InternalResponseHeaders).
		Then()

	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.Out.URL.Scheme = targetUrl.Scheme
			req.Out.URL.Host = targetUrl.Host
			req.Out.Host = targetUrl.Host
			
			req.SetXForwarded()
		},
		ErrorLog: logger.StdError(),
		ModifyResponse: modifyResponse,
	}

	return Proxy{Use: func() *httputil.ReverseProxy {return proxy}}
}

func (proxy Proxy) Init() http.Handler {
	handler := m.NewChain().
		Use(m.InternalInsecureAPI).
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
		Use(m.InternalSecureAPI).
		Then(proxy.Use())

	return handler
}
