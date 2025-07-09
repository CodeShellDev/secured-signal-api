package main

import (
	"net/http"
	"net/http/httputil"
	"os"

	proxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	env "github.com/codeshelldev/secured-signal-api/utils/env"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var initHandler *httputil.ReverseProxy

var ENV env.ENV_

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	env.Load()

	ENV = env.ENV

	initHandler = proxy.Create(ENV.API_URL)

	temp_m3 := TemplateMiddleware{
		Next:      initHandler,
		Variables: ENV.VARIABLES,
	}

	endp_m2 := EndpointsMiddleware{
		Next:             temp_m3.Use(),
		BlockedEndpoints: ENV.BLOCKED_ENDPOINTS,
	}

	auth_m1 := AuthMiddleware{
		Next:  endp_m2.Use(),
		Token: ENV.API_TOKEN,
	}

	log.Info("Initialized Proxy Handler")

	addr := "0.0.0.0:" + ENV.PORT

	log.Info("Server Listening on ", addr)

	http.ListenAndServe(addr, auth_m1.Use())
}
