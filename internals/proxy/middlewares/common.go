package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/utils/config"
)

type Context struct {
	Next  	http.Handler
}

type authType string

const (
	Bearer authType = "Bearer"
	Basic  authType = "Basic"
	Query  authType = "Query"
	None   authType = "None"
)

type contextKey string

const tokenKey contextKey = "token"

func getSettings(req *http.Request) config.SETTING_ {
	token := req.Context().Value(tokenKey).(string)

	settings := config.ENV.SETTINGS[token]

	if settings == nil {
		settings = config.ENV.SETTINGS["*"]
	}

	return *settings
}