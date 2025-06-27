package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	log "github.com/codeshelldev/secured-signal-api/utils"
)

type AuthType string

const (
	Bearer AuthType = "Bearer"
	Basic AuthType = "Basic"
	Query AuthType = "Query"
	None AuthType = "None"
)

func getAuthType(str string) AuthType {
	switch str {
	case "Bearer":
		return Bearer
	case "Basic":
		return Basic
	default:
		return None
	}
}

func renderTemplate(name string, tmplStr string, data any) (string, error) {
	tmpl, err := template.New(name).Parse(tmplStr)

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Info("Request:", req.Method, req.URL.Path)

		token := os.Getenv("API_TOKEN")
		user := "api"

		authHeader := req.Header.Get("Authorization")

		authQuery := req.URL.Query().Get("@authorization")

		var authType AuthType = None

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

func BlockedEndpointMiddleware(next http.Handler, BLOCKED_ENDPOINTS []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqPath := req.URL.Path

		if slices.Contains(BLOCKED_ENDPOINTS, reqPath) {
			log.Warn("User tried to access blocked endpoint: ", reqPath)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func TemplatingMiddleware(next http.Handler, VARIABLES map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)

			if err != nil {
				log.Error("Could not read Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			req.Body.Close()

			modifiedBody := string(bodyBytes)

			modifiedBody, _ = renderTemplate("json", modifiedBody, VARIABLES)

			modifiedBodyBytes := []byte(modifiedBody)

			if req.URL.RawQuery != "" {
				var modifiedBodyData map[string]interface{}

				err = json.Unmarshal(modifiedBodyBytes, &modifiedBodyData)

				if err != nil {
					log.Error("Could not decode Body: ", err.Error())
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}

				query, _ := renderTemplate("query", req.URL.RawQuery, VARIABLES)

				queryData, _ := url.ParseQuery(query)

				modifiedQuery := req.URL.Query()
				
				for key, value := range queryData {
					keyWithoutPrefix, found := strings.CutPrefix(key, "@")
	
					if found {
						modifiedBodyData[keyWithoutPrefix] = value

						modifiedQuery.Del(key)
					}
				}

				req.URL.RawQuery = modifiedQuery.Encode()

				modifiedBodyBytes, err = json.Marshal(modifiedBodyData)

				if err != nil {
					log.Error("Could not encode Body: ", err.Error())
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}

				log.Debug("Applied Query Templating: ", query)
			}

			log.Debug("Applied Body Templating")

			req.Body = io.NopCloser(bytes.NewReader(modifiedBodyBytes))

			req.ContentLength = int64(len(string(modifiedBodyBytes)))
			req.Header.Set("Content-Length", strconv.Itoa(len(modifiedBody)))
		}

		reqPath := req.URL.Path
		reqPath, _ = url.PathUnescape(reqPath)

		modifiedReqPath, _ := renderTemplate("path", reqPath, VARIABLES)

		req.URL.Path = modifiedReqPath

		next.ServeHTTP(w, req)
	})
}

func Create(targetUrl string) *httputil.ReverseProxy {
	url, _ := url.Parse(targetUrl)

	proxy := httputil.NewSingleHostReverseProxy(url)

	return proxy
}