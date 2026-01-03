package middlewares

import (
	"context"
	"net/http"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
)

type Context struct {
	Next http.Handler
}

type contextKey string

func setContext(req *http.Request, key, value any) {
	ctx := context.WithValue(req.Context(), tokenKey, value)
	req = req.WithContext(ctx)
}

func getContext[T any](req *http.Request, key any) T {
	value, ok := req.Context().Value(key).(T)

	if !ok {
		var zero T
		return zero
	}

	return value
}

func getLogger(req *http.Request) *logger.Logger {
	return getContext[*logger.Logger](req, loggerKey)
}

func getToken(req *http.Request) string {
	return getContext[string](req, tokenKey)
}

func getConfigByReq(req *http.Request) *structure.CONFIG {
	return getConfig(getToken(req))
}

func getConfig(token string) *structure.CONFIG {
	conf, exists := config.ENV.CONFIGS[token]

	if !exists || conf == nil {
		conf = config.DEFAULT
	}

	return conf
}
