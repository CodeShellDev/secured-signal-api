package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
	"github.com/codeshelldev/secured-signal-api/utils/requestkeys"
)

var Auth Middleware = Middleware{
	Name: "Auth",
	Use: authHandler,
}

type AuthAttempt struct {
	Error error
	Method *AuthMethod
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

		if strings.EqualFold(headerParts[0], "bearer") {
			req.Header.Del("Authorization")

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

		if strings.EqualFold(headerParts[0], "basic") {
			req.Header.Del("Authorization")

			base64Bytes, err := base64.StdEncoding.DecodeString(headerParts[1])

			if err != nil {
				logger.Error("Could not decode basic auth payload: ", err.Error())
				return "", errors.New("invalid base64 in basic auth")
			}

			parts := strings.SplitN(string(base64Bytes), ":", 2)

			if len(parts) != 2 {
				return "", errors.New("basic auth must be user:password")
			}

			user, password := parts[0], parts[1]

			if strings.EqualFold(user, "api") && isValidToken(tokens, password) {
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

			body.UpdateReq(req)

			return auth, nil
		}

		return "", errors.New("invalid body token")
	},
}

var QueryAuth = AuthMethod{
	Name: "Query",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		const authQuery = "auth"

		auth := req.URL.Query().Get(requestkeys.BodyPrefix + authQuery)

		if strings.TrimSpace(auth) == "" {
			return "", nil
		}

		if isValidToken(tokens, auth) {
			query := req.URL.Query()

			query.Del(requestkeys.BodyPrefix + authQuery)

			req.URL.RawQuery = query.Encode()

			return auth, nil
		}

		return "", errors.New("invalid query token")
	},
}

var PathAuth = AuthMethod{
	Name: "Path",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		const authPath = "auth"

		parts := strings.Split(req.URL.Path, "/")

		if len(parts) <= 1 {
			return "", nil
		}

		parts = parts[1:]

		unescaped, err := url.PathUnescape(parts[0])

		if err != nil {
			return "", nil
		}

		auth, exists := strings.CutPrefix(unescaped, "@" + authPath + "=")

		if !exists {
			return "", nil
		}

		req.URL.Path = "/" + strings.Join(parts[1:], "/")

		if isValidToken(tokens, auth) {
			return auth, nil
		}

		return "", errors.New("invalid path token")
	},
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

func (chain *AuthChain) Eval(w http.ResponseWriter, req *http.Request, tokens []string) (AuthMethod, string, error) {
	var err error
	var token string

	for _, method := range chain.methods {
		token, err = method.Authenticate(w, req, tokens)

		if err != nil {
			return method, "", err
		}

		if token != "" {
			return method, token, nil
		}
	}

	return AuthMethod{}, "", errors.New("no auth provided")
}

func authHandler(next http.Handler) http.Handler {
	var authChain = NewAuthChain().
		Use(BearerAuth).
		Use(BasicAuth).
		Use(BodyAuth).
		Use(QueryAuth).
		Use(PathAuth)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokens := config.ENV.TOKENS

		if tokens == nil {
			tokens = []string{}
		}

		if config.ENV.INSECURE || len(tokens) <= 0 {
			next.ServeHTTP(w, req)
			return
		}

		method, token, err := authChain.Eval(w, req, tokens)

		if token == "" {
			if method.Name != "" {
				req = SetContext(req, AuthAttemptKey, AuthAttempt{
					Error: err,
					Method: &method,
				})
			} else {
				req = SetContext(req, AuthAttemptKey, AuthAttempt{
					Error: err,
				})
			}

			req = SetContext(req, IsAuthKey, false)
		} else {
			conf := GetConfigWithoutDefault(token)

			allowedMethods := conf.API.AUTH.METHODS.OptOrEmpty(config.DEFAULT.API.AUTH.METHODS)

			if isAuthMethodAllowed(method, token, allowedMethods, conf.API.AUTH.TOKENS) {
				req = SetContext(req, IsAuthKey, true)
				req = SetContext(req, TokenKey, token)
			} else {
				req = SetContext(req, AuthAttemptKey, AuthAttempt{
					Error: errors.New("disabled auth method"),
					Method: &method,
				})

				req = SetContext(req, IsAuthKey, false)
			}
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
		isAuthenticated := GetContext[bool](req, IsAuthKey)

		logger.Dev(req.Context())

		if !isAuthenticated {
			attempt := GetContext[AuthAttempt](req, AuthAttemptKey)

			if attempt.Method != nil {
				logger.Warn("Client failed ", attempt.Method.Name, " auth: ", attempt.Error.Error())
			} else {
				logger.Dev("IsAuth: ", isAuthenticated)
				logger.Warn("Client failed to authenticate: ", attempt.Error.Error())
			}

			onUnauthorized(w)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func onUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Login Required\", Bearer realm=\"Access Token Required\"")

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func isValidToken(tokens []string, match string) bool {
	return slices.Contains(tokens, match)
}

type AuthToken struct {
	Token		string
	Methods		[]string
}

func isAuthMethodAllowed(method AuthMethod, token string, defaultMethods []string, tokenOverwrites []structure.Token) bool {
	if len(defaultMethods) == 0 && len(tokenOverwrites) == 0 {
		// default: allow all
		return true
	}

	for _, t := range tokenOverwrites {
		if slices.Contains(t.Set, token) {
			return slices.ContainsFunc(t.Methods, func(try string) bool {
				return strings.EqualFold(try, method.Name)
			})
		}
	}

	return slices.ContainsFunc(defaultMethods, func(try string) bool {
		return strings.EqualFold(try, method.Name)
	})
}