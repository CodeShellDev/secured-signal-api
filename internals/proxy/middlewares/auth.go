package middlewares

import (
	"context"
	"encoding/base64"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	log "github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config"
)

var Auth Middleware = Middleware{
	Name: "Auth",
	Use: authHandler,
}

type AuthMethod struct {
	Name string
	Authenticate func(req *http.Request, tokens []string) (bool, error)
}

var BearerAuth = AuthMethod {
	Name: "Bearer",
	Authenticate: func(req *http.Request, tokens []string) (bool, error) {
		header := req.Header.Get("Authorization")

		headerParts := strings.SplitN(header, " ", 2)

		if len(headerParts) != 2 {
			return false, nil
		}

		if strings.ToLower(headerParts[0]) == "bearer" {
			if isValidToken(tokens, headerParts[1]) {
				return true, nil
			}

			return false, errors.New("invalid Bearer token")
		}

		return false, nil
	},
}

var BasicAuth = AuthMethod {
	Name: "Basic",
	Authenticate: func(req *http.Request, tokens []string) (bool, error) {
		header := req.Header.Get("Authorization")

		if strings.TrimSpace(header) == "" {
			return false, nil
		}

		headerParts := strings.SplitN(header, " ", 2)

		if len(headerParts) != 2 {
			return false, nil
		}

		if strings.ToLower(headerParts[0]) == "basic" {
			base64Bytes, err := base64.StdEncoding.DecodeString(headerParts[1])

			if err != nil {
				log.Error("Could not decode Basic auth payload: ", err.Error())
				return false, errors.New("invalid base64 in Basic auth")
			}

			parts := strings.SplitN(string(base64Bytes), ":", 2)

			if len(parts) != 2 {
				return false, errors.New("Basic auth must be user:password")
			}

			user, password := parts[0], parts[1]

			if strings.ToLower(user) == "api" && isValidToken(tokens, password) {
				return true, nil
			}

			return false, errors.New("invalid user:password")
		}

		return false, nil
	},
}

var QueryAuth = AuthMethod {
	Name: "Query",
	Authenticate: func(req *http.Request, tokens []string) (bool, error) {
		const authQuery = "@authorization"

		auth := req.URL.Query().Get(authQuery)

		if strings.TrimSpace(auth) == "" {
			return false, nil
		}

		if isValidToken(tokens, auth) {
			query := req.URL.Query()

			query.Del(authQuery)

			req.URL.RawQuery = query.Encode()

			return true, nil
		}

		return false, errors.New("invalid Query token")
	},
}

var PathAuth = AuthMethod {
	Name: "Path",
	Authenticate: func(req *http.Request, tokens []string) (bool, error) {
		parts := strings.Split(req.URL.Path, "/")

		if len(parts) == 0 {
			return false, nil
		}

		unescaped, err := url.PathUnescape(parts[1])

		if err != nil {
			return false, nil
		}

		auth, exists := strings.CutPrefix(unescaped, "auth=")

		if !exists {
			return false, nil
		}

		if isValidToken(tokens, auth) {
			return true, nil
		}

		return false, errors.New("invalid Path token")
	},
}

func authHandler(next http.Handler) http.Handler {
	tokenKeys := maps.Keys(config.ENV.CONFIGS)
	tokens := slices.Collect(tokenKeys)

	if tokens == nil {
		tokens = []string{}
	}

	var authChain = NewAuthChain().
		Use(BearerAuth).
		Use(BasicAuth).
		Use(QueryAuth).
		Use(PathAuth)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(tokens) <= 0 {
			next.ServeHTTP(w, req)
			return
		}

		var authToken string

		success, _ := authChain.Eval(req, tokens)

		if !success {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Login Required\", Bearer realm=\"Access Token Required\"")

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), tokenKey, authToken)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func isValidToken(tokens []string, match string) bool {
	return slices.Contains(tokens, match)
}

type AuthChain struct {
    methods []AuthMethod
}

func NewAuthChain() *AuthChain {
    return &AuthChain{}
}

func (chain *AuthChain) Use(method AuthMethod) *AuthChain {
    chain.methods = append(chain.methods, method)

	logger.Debug("Registered ", method.Name, " auth")

    return chain
}

func (chain *AuthChain) Eval(req *http.Request, tokens []string) (bool, error) {
	var err error
	var success bool

	for _, method := range chain.methods {
		success, err = method.Authenticate(req, tokens)

		if err != nil {
			logger.Warn("User failed ", method.Name, " auth: ", err.Error())
		}

		if success {
			return success, nil
		}
	}

	return false, err
}