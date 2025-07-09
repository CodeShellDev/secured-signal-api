package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	query "github.com/codeshelldev/secured-signal-api/utils/query"
)

type TemplateMiddleware struct {
	Next      http.Handler
	Variables map[string]interface{}
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
				for i, tmplStr := range matches {

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

func (data TemplateMiddleware) Use() http.Handler {
	next := data.Next
	VARIABLES := data.Variables

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
				decodedQuery, _ := url.QueryUnescape(req.URL.RawQuery)

				log.Debug("Decoded Query: ", decodedQuery)

				templatedQuery, _ := renderTemplate("query", decodedQuery, VARIABLES)

				modifiedQuery := req.URL.Query()

				queryData := query.ParseRawQuery(templatedQuery)

				for key, value := range queryData {
					keyWithoutPrefix, found := strings.CutPrefix(key, "@")

					if found {
						modifiedBodyData[keyWithoutPrefix] = query.ParseTypedQuery(value)

						modifiedQuery.Del(key)
					}
				}

				req.URL.RawQuery = modifiedQuery.Encode()

				log.Debug("Applied Query Templating: ", templatedQuery)
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
