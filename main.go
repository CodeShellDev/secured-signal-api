package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"os"

	proxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	log "github.com/codeshelldev/secured-signal-api/utils"
)

var handler *httputil.ReverseProxy

var VARIABLES map[string]string = map[string]string{
	"RECIPIENTS": os.Getenv("DEFAULT_RECIPIENTS"),
	"NUMBER": os.Getenv("SENDER"),
}

var BLOCKED_ENDPOINTS []string = []string{
    "/v1/about",
    "/v1/configuration",
    "/v1/devices",
    "/v1/register",
    "/v1/unregister",
    "/v1/qrcodelink",
    "/v1/accounts",
    "/v1/contacts",
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	port := os.Getenv("PORT")
	signalUrl := os.Getenv("SIGNAL_API_URL")

	blockedEndpointJSON := os.Getenv("BLOCKED_ENDPOINTS")
	variablesJSON := os.Getenv("VARIABLES")

	log.Info("Loaded Environment Variables")

	if blockedEndpointJSON != "" {
		var blockedEndpoints []string

		err := json.Unmarshal([]byte(blockedEndpointJSON), &blockedEndpoints)

		if err != nil {
			log.Error("Could not decode Blocked Endpoints: ", blockedEndpointJSON)
		}

		BLOCKED_ENDPOINTS = blockedEndpoints
	}

	if variablesJSON != "" {
		var variables map[string]string

		err := json.Unmarshal([]byte(variablesJSON), &variables)

		if err != nil {
			log.Error("Could not decode Variables ", variablesJSON)
		}

		VARIABLES = variables
	}

	handler = proxy.Create(signalUrl)

	finalHandler := proxy.AuthMiddleware(
		proxy.BlockedEndpointMiddleware(
			proxy.TemplatingMiddleware(handler,
				VARIABLES ),
				
		BLOCKED_ENDPOINTS ),
	)

	log.Info("Initialized Proxy Handler")

	addr := "0.0.0.0:" + port

	log.Info("Server Listening on ", addr)

	http.ListenAndServe(addr, finalHandler)
}