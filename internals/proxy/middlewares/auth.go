package middlewares

import (
	"context"
	"encoding/base64"
	"net/http"
	"slices"
	"strings"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var Auth Middleware = Middleware{
	Name: "Auth",
	Use: authHandler,
}

func authHandler(next http.Handler) http.Handler {
	tokens := config.ENV.API_TOKENS

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(tokens) <= 0 {
			next.ServeHTTP(w, req)
			return
		}

		authHeader := req.Header.Get("Authorization")

		authQuery := req.URL.Query().Get("@authorization")

		var authType authType = None

		var authToken string

		success := false

		if authHeader != "" {
			authBody := strings.Split(authHeader, " ")

			authType = getAuthType(authBody[0])
			authToken = authBody[1]

			switch authType {
			case Bearer:
				if isValidToken(tokens, authToken) {
					success = true
				}

			case Basic:
				basicAuthBody, err := base64.StdEncoding.DecodeString(authToken)

				if err != nil {
					log.Error("Could not decode Basic Auth Payload: ", err.Error())
				}

				basicAuth := string(basicAuthBody)
				basicAuthParams := strings.Split(basicAuth, ":")

				user := "api"

				if basicAuthParams[0] == user && isValidToken(tokens, basicAuthParams[1]) {
					success = true
				}
			}

		} else if authQuery != "" {
			authType = Query

			authToken = strings.TrimSpace(authQuery)

			if isValidToken(tokens, authToken) {
				success = true

				modifiedQuery := req.URL.Query()

				modifiedQuery.Del("@authorization")

				req.URL.RawQuery = modifiedQuery.Encode()
			}
		}

		if !success {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Login Required\", Bearer realm=\"Access Token Required\"")

			log.Warn("User failed ", string(authType), " Auth")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), tokenKey, authToken)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func getAuthType(str string) authType {
	switch str {
	case "Bearer":
		return Bearer
	case "Basic":
		return Basic
	default:
		return None
	}
}

func isValidToken(tokens []string, match string) bool {
	return slices.Contains(tokens, match)
}