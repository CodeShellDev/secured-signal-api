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
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type TemplateMiddleware struct {
	Next      http.Handler
	Variables map[string]interface{}
}

func (data TemplateMiddleware) Use() http.Handler {
	next := data.Next
	VARIABLES := data.Variables

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
		}

		bodyData := map[string]interface{}{}

		var modifiedBody bool

		if !body.Empty {
			var modified bool

			bodyData, modified = templateJSON(body.Data, VARIABLES)

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.RawQuery != "" {
			var modified bool

			req.URL.RawQuery, bodyData, modified = templateQuery(bodyData, req.URL, VARIABLES)

			if modified {
				modifiedBody = true
			}
		}

		if modifiedBody {
			modifiedBody, err := request.CreateBody(bodyData)

			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			body = modifiedBody

			strData := body.ToString()

			log.Debug("Applied Body Templating: ", strData)

			req.ContentLength = int64(len(strData))
			req.Header.Set("Content-Length", strconv.Itoa(len(strData)))
		}

		req.Body = io.NopCloser(bytes.NewReader(body.Raw))

		req.URL.Path, _ = templatePath(req.URL, VARIABLES)

		next.ServeHTTP(w, req)
	})
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

func templateJSON(data map[string]interface{}, variables map[string]interface{}) (map[string]interface{}, bool) {
	var modified bool

	for k, v := range data {
		str, ok := v.(string)

		if ok {
			re, err := regexp.Compile(`{{\s*\.([A-Za-z_][A-Za-z0-9_]*)\s*}}`)

			if err != nil {
				log.Error("Error while Compiling Regex: ", err.Error())
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

				modified = true
			} else if len(matches) == 1 {
				tmplKey := matches[0][1]

				data[k] = variables[tmplKey]

				modified = true
			}
		}
	}

	return data, modified
}

func templatePath(reqUrl *url.URL, VARIABLES interface{}) (string, bool) {
	var modified bool

	reqPath, err := url.PathUnescape(reqUrl.Path)

	if err != nil {
		log.Error("Error while Escaping Path: ", err.Error())
		return reqUrl.Path, modified
	}

	reqPath, err = renderTemplate("path", reqPath, VARIABLES)

	if err != nil {
		log.Error("Could not Template Path: ", err.Error())
		return reqUrl.Path, modified
	}

	if reqUrl.Path != reqPath {
		log.Debug("Applied Path Templating: ", reqPath)

		modified = true
	}

	return reqPath, modified
}

func templateQuery(data map[string]interface{}, reqUrl *url.URL, VARIABLES interface{}) (string, map[string]interface{}, bool) {
	var modified bool

	decodedQuery, _ := url.QueryUnescape(reqUrl.RawQuery)

	log.Debug("Decoded Query: ", decodedQuery)

	templatedQuery, _ := renderTemplate("query", decodedQuery, VARIABLES)

	modifiedQuery := reqUrl.Query()

	queryData := query.ParseRawQuery(templatedQuery)

	for key, value := range queryData {
		keyWithoutPrefix, found := strings.CutPrefix(key, "@")

		if found {
			data[keyWithoutPrefix] = query.ParseTypedQuery(value)

			modifiedQuery.Del(key)
		}
	}

	reqRawQuery := modifiedQuery.Encode()

	if reqUrl.Query().Encode() != reqRawQuery {
		log.Debug("Applied Query Templating: ", templatedQuery)

		modified = true
	}

	return reqRawQuery, data, modified
}