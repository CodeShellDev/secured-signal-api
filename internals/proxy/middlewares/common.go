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

func getConfigByReq(req *http.Request) *structure.CONFIG {
	token := req.Context().Value(tokenKey).(string)

	return getConfig(token)
}

func getConfig(token string) *structure.CONFIG {
	conf, exists := config.ENV.CONFIGS[token]

	if !exists || conf == nil {
		conf = config.DEFAULT
	}

	return conf
}
