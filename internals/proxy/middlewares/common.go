package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/utils/config"
	"github.com/codeshelldev/secured-signal-api/utils/config/structure"
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

func getSettingsByReq(req *http.Request) *structure.SETTING_ {
	token, ok := req.Context().Value(tokenKey).(string)

	if !ok {
		token = "*"
	}

	return getSettings(token)
}

func getSettings(token string) *structure.SETTING_ {
	settings, exists := config.ENV.SETTINGS[token]

	if !exists || settings == nil {
		settings = config.ENV.SETTINGS["*"]
	}

	return settings
}
