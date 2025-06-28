package main

import (
	"net/http"
	"net/http/httputil"
	"os"

	proxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	env "github.com/codeshelldev/secured-signal-api/utils/env"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var handler *httputil.ReverseProxy

var ENV env.ENV_

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	env.Load()

	ENV = env.ENV

	handler = proxy.Create(ENV.API_URL)

	finalHandler := proxy.AuthMiddleware(
		proxy.BlockedEndpointMiddleware(
			proxy.TemplatingMiddleware(handler,
				ENV.VARIABLES ),
		ENV.BLOCKED_ENDPOINTS ),
	ENV.API_TOKEN )

	log.Info("Initialized Proxy Handler")

	addr := "0.0.0.0:" + ENV.PORT

	log.Info("Server Listening on ", addr)

	http.ListenAndServe(addr, finalHandler)
}