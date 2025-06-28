package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type AuthType string

const (
	Bearer AuthType = "Bearer"
	Basic AuthType = "Basic"
	Query AuthType = "Query"
	None AuthType = "None"
)

func parseTypedQuery(values []string) interface{} {
	var result interface{}

	raw := values[0]

	intValue, err := strconv.Atoi(raw)

	if strings.Contains(raw, ",") {
		parts := strings.Split(raw, ",")
		var list []interface{}
		for _, part := range parts {
			if intVal, err := strconv.Atoi(part); err == nil {
				list = append(list, intVal)
			} else {
				list = append(list, part)
			}
		}
		result = list
	} else if err == nil {
		result = intValue
	} else {
		result = raw
	}

	return result
}

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

	// TODO: Escape Arrays inside of strings "{{ .ARRAY }}" => [ 1, 2, 3 ]

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

func templateJSON(data map[string]interface{}, variables map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		str, ok := v.(string)

		if ok {
			re, err := regexp.Compile(`{{\s*\.([A-Za-z_][A-Za-z0-9_]*)\s*}}`)

			if err != nil {
				log.Error("Encountered Error while Compiling Regex: ", err.Error())
			}

			matches := re.FindAllStringSubmatch(str, -1)

			if len(matches) > 1 {
				for i, tmplStr := range(matches) {
					
					tmplKey := matches[i][1]

					variable, err := json.Marshal(variables[tmplKey])

					if err != nil {
						log.Error("Could not decode JSON: ", err.Error())
						break
					}

					data[k] = strings.ReplaceAll(str, string(variable), tmplStr[0])
				}
			} else if len(matches) == 1 {
				tmplKey := matches[0][1]

				data[k] = variables[tmplKey]
			}
		}
	}

	return data
}

func AuthMiddleware(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if token == "" {
			next.ServeHTTP(w, req)
			return
		}

		log.Info("Request:", req.Method, req.URL.Path)

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

func TemplatingMiddleware(next http.Handler, VARIABLES map[string]interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)

			if err != nil {
				log.Error("Could not read Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			req.Body.Close()

			var modifiedBodyData map[string]interface{}

			err = json.Unmarshal(bodyBytes, &modifiedBodyData)

			if err != nil {
				log.Error("Could not decode Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			modifiedBodyData = templateJSON(modifiedBodyData, VARIABLES)

			if req.URL.RawQuery != "" {
				query, _ := renderTemplate("query", req.URL.RawQuery, VARIABLES)

				modifiedQuery := req.URL.Query()

				queryData, _ := url.ParseQuery(query)

				for key, value := range queryData {
					keyWithoutPrefix, found := strings.CutPrefix(key, "@")
	
					if found {
						modifiedBodyData[keyWithoutPrefix] = parseTypedQuery(value)

						modifiedQuery.Del(key)
					}
				}

				req.URL.RawQuery = modifiedQuery.Encode()

				log.Debug("Applied Query Templating: ", query)
			}

			modifiedBodyBytes, err := json.Marshal(modifiedBodyData)

			if err != nil {
				log.Error("Could not encode Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			modifiedBody := string(modifiedBodyBytes)

			log.Debug("Applied Body Templating: ", modifiedBody)

			req.Body = io.NopCloser(bytes.NewReader(modifiedBodyBytes))

			req.ContentLength = int64(len(modifiedBody))
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