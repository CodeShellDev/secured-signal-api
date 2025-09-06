package main

import (
	"net/http"
	"net/http/httputil"
	"os"

	proxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
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

	body_m4 := middlewares.BodyMiddleware{
		Next: initHandler,
		MessageAliases: ENV.MESSAGE_ALIASES,
	}

	temp_m3 := middlewares.TemplateMiddleware{
		Next:      body_m4.Use(),
		Variables: ENV.VARIABLES,
	}

	endp_m2 := middlewares.EndpointsMiddleware{
		Next:             temp_m3.Use(),
		BlockedEndpoints: ENV.BLOCKED_ENDPOINTS,
	}

	auth_m1 := middlewares.AuthMiddleware{
		Next:   endp_m2.Use(),
		Tokens: ENV.API_TOKENS,
	}

	log_m0 := middlewares.LogMiddleware{
		Next: auth_m1.Use(),
	}

	log.Info("Initialized Proxy Handler")

	addr := "0.0.0.0:" + ENV.PORT

	log.Info("Server Listening on ", addr)

	http.ListenAndServe(addr, log_m0.Use())
}
