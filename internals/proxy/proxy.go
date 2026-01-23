package proxy

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/codeshelldev/gotl/pkg/ioutils"
	"github.com/codeshelldev/gotl/pkg/logger"
	m "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
)

type Proxy struct {
	Use func() *httputil.ReverseProxy
}

func Create(targetUrl string) Proxy {
	if strings.TrimSpace(targetUrl) == "" {
		logger.Fatal("Missing API URL")
		return Proxy{Use: func() *httputil.ReverseProxy {return nil}}
	}

	url, err := url.Parse(targetUrl)

	if err != nil {
		logger.Fatal("Invalid API URL: ", targetUrl)
		return Proxy{Use: func() *httputil.ReverseProxy {return nil}}
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(res *http.Response) error {
		res.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private, proxy-revalidate")
		res.Header.Set("Pragma", "no-cache")
		res.Header.Set("Expires", "0")
		res.Header.Set("Vary", "*")
		res.Header.Set("Referrer-Policy", "no-referrer")

		return nil
	}

	w := &ioutils.InterceptWriter{
		Writer: &bytes.Buffer{},
		Hook: func(bytes []byte) {
			msg := string(bytes)
			msg = strings.TrimSuffix(msg, "\n") 

			logger.Error(msg)
		},
	}

	proxy.ErrorLog = log.New(w, "", 0)

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
		Use(m.Server).
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
