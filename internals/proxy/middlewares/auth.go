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
	"github.com/codeshelldev/secured-signal-api/utils/deprecation"
)

var Auth Middleware = Middleware{
	Name: "Auth",
	Use: authHandler,
}

const tokenKey contextKey = "token"
const isAuthKey contextKey = "isAuthenticated"

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

		if strings.ToLower(headerParts[0]) == "basic" {
			req.Header.Del("Authorization")

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

		return "", errors.New("invalid Body token")
	},
}

var QueryAuth = AuthMethod{
	Name: "Query",
	Authenticate: func(w http.ResponseWriter, req *http.Request, tokens []string) (string, error) {
		const authQuery = "auth"

		auth := req.URL.Query().Get("@" + authQuery)

		// BREAKING @authorization Query
		const oldAuthQuery = "authorization"

		if req.URL.Query().Has("@" + oldAuthQuery) {
			fullURL, _ := request.ParseReqURL(req)
			urlWithNewAuthQuery := strings.Replace(fullURL.String(), "@" + oldAuthQuery, "@{s,fg=bright_red}" + oldAuthQuery + "{/}{b,fg=green}" + authQuery + "{/}", 1)

			deprecation.Error(req.URL.String(), deprecation.DeprecationMessage{
				Using: "{b,i,bg=red}`@authorization`{/} in the query",
				Message: "{b,fg=red}`/?@{s}authorization{/}`{/} has been renamed to {b,fg=green}`/?@auth`{}",
				Fix: "\nChange the {b}url{/} to:\n`" + urlWithNewAuthQuery + "`",
			})
		}

		if strings.TrimSpace(auth) == "" {
			return "", nil
		}

		if isValidToken(tokens, auth) {
			query := req.URL.Query()

			query.Del("@" + authQuery)

			req.URL.RawQuery = query.Encode()

			return auth, nil
		}

		return "", errors.New("invalid Query token")
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

		return "", errors.New("invalid Path token")
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
			logger.Warn("Client failed ", method.Name, " auth: ", err.Error())
			return AuthMethod{}, "", err
		}

		if token != "" {
			return method, token, nil
		}
	}

	logger.Warn("Client failed to provide any auth")

	return AuthMethod{}, "", err
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

		method, token, _ := authChain.Eval(w, req, tokens)

		if token == "" {
			onUnauthorized(w)

			req = setContext(req, isAuthKey, false)
		} else {
			conf := getConfigWithoutDefault(token)

			allowedMethods := conf.API.AUTH.METHODS

			if allowedMethods == nil {
				allowedMethods = getConfig("").API.AUTH.METHODS
			}

			if isAuthMethodAllowed(method, token, conf.API.TOKENS, allowedMethods, conf.API.AUTH.TOKENS) {
				req = setContext(req, isAuthKey, true)
				req = setContext(req, tokenKey, token)
			} else {
				// BREAKING Query & Path auth disabled (default)
				if (method.Name == "Path" || method.Name == "Query") && conf.API.AUTH.METHODS == nil {
					deprecation.Error(method.Name, deprecation.DeprecationMessage{
						Message: "{b}Query{/} and {b}Path{/} auth are {u}disabled{/} by default\nTo be able to use them they must first be enabled",
						Fix: "\n{b}Add{/} {b,fg=green}`" + strings.ToLower(method.Name) + "`{/} to {i}`api.auth.methods`{/}:" + 
							"\napi.auth.methods: [" + strings.Join(append(allowedMethods, "{b,fg=green}" + strings.ToLower(method.Name) + "{/}"), ", ") + "]",
						Note: "\n{i}Let us know what you think about this change at\n{i}{u,fg=blue}https://github.com/CodeShellDev/secured-signal-api/discussions/221{/}{/}",
					})
				}

				logger.Warn("Client tried using disabled auth method: ", method.Name)

				onUnauthorized(w)

				req = setContext(req, isAuthKey, false)
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
		isAuthenticated := getContext[bool](req, isAuthKey)

		if !isAuthenticated {
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

func getTokenMethodMap(rawTokens []string, defaultMethods []string, tokenMethodSet []structure.Token) map[string][]string {
	tokenMethodMap := map[string][]string{}

	for _, token := range rawTokens {
		tokenMethodMap[token] = defaultMethods
	}

	for _, set := range tokenMethodSet {
		for _, token := range set.Set {
			tokenMethodMap[token] = set.Methods
		}
	}

	return tokenMethodMap
}

func isAuthMethodAllowed(method AuthMethod, token string, rawTokens []string, defaultMethods []string, tokenMethodSet []structure.Token) bool {
	if (len(defaultMethods) == 0 || defaultMethods == nil) && (len(tokenMethodSet) == 0 || tokenMethodSet == nil) {
		// default: allow all
		return true
	}

	tokenMethodMap := getTokenMethodMap(rawTokens, defaultMethods, tokenMethodSet)

	return slices.ContainsFunc(tokenMethodMap[token], func(try string) bool {
		return strings.EqualFold(try, method.Name)
	})
}