package middleware

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var InternalResponseHeaders = ResponseMiddleware{
	Name: "_Response_Headers",
	Use: headersHandler,
}

func headersHandler(res *http.Response) error {
	conf := GetConfigByReq(res.Request)

	resHeaders := conf.SETTINGS.HTTP.RESPONSE_HEADERS.OptOrEmpty(config.DEFAULT.SETTINGS.HTTP.RESPONSE_HEADERS)

	res.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private, proxy-revalidate")
	res.Header.Set("Pragma", "no-cache")
	res.Header.Set("Expires", "0")
	res.Header.Set("Vary", "*")
	res.Header.Set("Referrer-Policy", "no-referrer")

	if len(resHeaders) != 0 {
		for k, v := range resHeaders {
			res.Header.Set(k, v)
		}
	}

	return nil
}