package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
)

type Context struct {
	Next http.Handler
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

func getSettingsByReq(req *http.Request) *structure.SETTINGS {
	token, ok := req.Context().Value(tokenKey).(string)

	if !ok {
		token = "*"
	}

	return getSettings(token)
}

func getSettings(token string) *structure.SETTINGS {
	settings, exists := config.ENV.SETTINGS[token]

	if !exists || settings == nil {
		settings = config.ENV.SETTINGS["*"]
	}

	return settings
}
