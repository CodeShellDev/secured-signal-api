package middlewares

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type AuthMiddleware struct {
	Next  http.Handler
	Token string
}

type authType string

const (
	Bearer authType = "Bearer"
	Basic  authType = "Basic"
	Query  authType = "Query"
	None   authType = "None"
)

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

func (data AuthMiddleware) Use() http.Handler {
	next := data.Next
	token := data.Token

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if token == "" {
			next.ServeHTTP(w, req)
			return
		}

		log.Info("Request:", req.Method, req.URL.Path)

		authHeader := req.Header.Get("Authorization")

		authQuery := req.URL.Query().Get("@authorization")

		var authType authType = None

		success := false

		if authHeader != "" {
			authBody := strings.Split(authHeader, " ")

			authType = getAuthType(authBody[0])
			authToken := authBody[1]

			switch authType {
			case Bearer:
				if authToken == token {
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

				if basicAuthParams[0] == user && basicAuthParams[1] == token {
					success = true
				}
			}

		} else if authQuery != "" {
			authType = Query

			authToken, _ := url.QueryUnescape(authQuery)

			if authToken == token {
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

		next.ServeHTTP(w, req)
	})
}
