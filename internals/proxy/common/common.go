package common

import (
	"context"
	"net/http"
	"net/url"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
)

func SetContext(req *http.Request, key, value any) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

func GetContext[T any](req *http.Request, key any) T {
	value, ok := req.Context().Value(key).(T)

	if !ok {
		var zero T
		return zero
	}

	return value
}

func GetLogger(req *http.Request) *logger.Logger {
	return GetContext[*logger.Logger](req, LoggerKey)
}

func GetToken(req *http.Request) string {
	return GetContext[string](req, TokenKey)
}

func GetConfigByReq(req *http.Request) *structure.CONFIG {
	return GetConfig(GetToken(req))
}

func GetConfigWithoutDefaultByReq(req *http.Request) *structure.CONFIG {
	return GetConfigWithoutDefault(GetToken(req))
}

func GetConfigWithoutDefault(token string) *structure.CONFIG {
	conf, exists := config.ENV.CONFIGS[token]

	if !exists {
		return nil
	}

	return conf
}

func GetConfig(token string) *structure.CONFIG {
	conf := GetConfigWithoutDefault(token)

	if conf == nil {
		conf = config.DEFAULT
	}

	return conf
}


func ChangeRequestDest(req *http.Request, newDest string) error {
	newURL, err := url.Parse(newDest)
	if err != nil {
		return err
	}

	req.URL = newURL
	req.Host = newURL.Host

	return nil
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	res := request.Body{
		Data: map[string]any{
			"error": msg,
		},
	}

	w.WriteHeader(status)
	res.Write(w)
}