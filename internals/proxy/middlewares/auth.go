package middlewares

import (
	"encoding/base64"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
)

var Auth Middleware = Middleware{
	Name: "Auth",
	Use: authHandler,
}

const tokenKey contextKey = "token"
const isAuthKey contextKey = "isAuthenticated"

func authHandler(next http.Handler) http.Handler {
	tokenKeys := maps.Keys(config.ENV.CONFIGS)
	tokens := slices.Collect(tokenKeys)

	if tokens == nil {
		tokens = []string{}
	}

	var authChain = NewAuthChain().
		Use(BearerAuth).
		Use(BasicAuth).
		Use(BodyAuth).
		Use(QueryAuth).
		Use(PathAuth)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if config.ENV.INSECURE || len(tokens) <= 0 {
			next.ServeHTTP(w, req)
			return
		}

		token, _ := authChain.Eval(w, req, tokens)

		if token == "" {
			onUnauthorized(w)

			req = setContext(req, isAuthKey, false)
		} else {
			req = setContext(req, isAuthKey, true)
			req = setContext(req, tokenKey, token)
		}

		next.ServeHTTP(w, req)
	})
}

var InternalAuthRequirement Middleware = Middleware{
	Name: "_Auth_Requirement",
	Use: authRequirementHandler,
}

func authRequirementHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		isAuthenticated := getContext[bool](req, isAuthKey)

		if !isAuthenticated {
			return
		}

		next.ServeHTTP(w, req)
	})
}

type AuthMethod struct {
	Name string
	Authenticate func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error)
}

var BearerAuth = AuthMethod{
	Name: "Bearer",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		header := req.Header.Get("Authorization")

		headerParts := strings.SplitN(header, " ", 2)

		if len(headerParts) != 2 {
			return "", nil
		}

		if strings.ToLower(headerParts[0]) == "bearer" {
			if isValidToken(tokens, headerParts[1]) {
				return headerParts[1], nil
			}

			return "", errors.New("invalid Bearer token")
		}

		return "", nil
	},
}

var BasicAuth = AuthMethod{
	Name: "Basic",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		header := req.Header.Get("Authorization")

		if strings.TrimSpace(header) == "" {
			return "", nil
		}

		headerParts := strings.SplitN(header, " ", 2)

		if len(headerParts) != 2 {
			return "", nil
		}

		if strings.ToLower(headerParts[0]) == "basic" {
			base64Bytes, err := base64.StdEncoding.DecodeString(headerParts[1])

			if err != nil {
				logger.Error("Could not decode Basic auth payload: ", err.Error())
				return "", errors.New("invalid base64 in Basic auth")
			}

			parts := strings.SplitN(string(base64Bytes), ":", 2)

			if len(parts) != 2 {
				return "", errors.New("Basic auth must be user:password")
			}

			user, password := parts[0], parts[1]

			if strings.ToLower(user) == "api" && isValidToken(tokens, password) {
				return password, nil
			}

			return "", errors.New("invalid user:password")
		}

		return "", nil
	},
}

var BodyAuth = AuthMethod{
	Name: "Body",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		const authField = "auth"

		body, err := request.GetReqBody(req)

		if err != nil {
			return "", nil
		}

		body.Write(req)

		if body.Empty {
			return "", nil
		}

		value, exists := body.Data[authField]

		if !exists {
			return "", nil
		}

		auth, ok := value.(string)

		if !ok {
			return "", nil
		}

		if isValidToken(tokens, auth) {
			delete(body.Data, authField)

			body.Write(req)

			return auth, nil
		}

		return "", errors.New("invalid Body token")
	},
}

var QueryAuth = AuthMethod{
	Name: "Query",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		const authQuery = "@authorization"

		auth := req.URL.Query().Get(authQuery)

		if strings.TrimSpace(auth) == "" {
			return "", nil
		}

		if isValidToken(tokens, auth) {
			query := req.URL.Query()

			query.Del(authQuery)

			req.URL.RawQuery = query.Encode()

			return auth, nil
		}

		return "", errors.New("invalid Query token")
	},
}

var PathAuth = AuthMethod{
	Name: "Path",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		parts := strings.Split(req.URL.Path, "/")

		if len(parts) == 0 {
			return "", nil
		}

		unescaped, err := url.PathUnescape(parts[1])

		if err != nil {
			return "", nil
		}

		auth, exists := strings.CutPrefix(unescaped, "auth=")

		if !exists {
			return "", nil
		}

		if isValidToken(tokens, auth) {
			return auth, nil
		}

		return "", errors.New("invalid Path token")
	},
}

func onUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Login Required\", Bearer realm=\"Access Token Required\"")

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

func (chain *AuthChain) Eval(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
	var err error
	var token string

	for _, method := range chain.methods {
		token, err = method.Authenticate(w, req, tokens)

		if err != nil {
			logger.Warn("User failed ", method.Name, " auth: ", err.Error())
		}

		if token != "" {
			return token, nil
		}
	}

	return "", err
}