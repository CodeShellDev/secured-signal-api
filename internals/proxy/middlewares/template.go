package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	query "github.com/codeshelldev/secured-signal-api/utils/query"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
	templating "github.com/codeshelldev/secured-signal-api/utils/templating"
)

type TemplateMiddleware struct {
	Next      http.Handler
}

func (data TemplateMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		variables := getSettingsByReq(req).VARIABLES

		if variables == nil {
			variables = getSettings("*").VARIABLES
		}

		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty {
			var modified bool

			bodyData, modified, err = TemplateBody(body.Data, variables)

			if err != nil {
				log.Error("Error Templating JSON: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.RawQuery != "" {
			var modified bool

			req.URL.RawQuery, bodyData, modified, err = TemplateQuery(req.URL, bodyData, variables)

			if err != nil {
				log.Error("Error Templating Query: ", err.Error())
			}

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

		if req.URL.Path != "" {
			var modified bool

			req.URL.Path, modified, err = TemplatePath(req.URL, variables)

			if err != nil {
				log.Error("Error Templating Path: ", err.Error())
			}

			if modified {
				log.Debug("Applied Path Templating: ", req.URL.Path)
			}
		}

		next.ServeHTTP(w, req)
	})
}

func TemplateBody(data map[string]any, VARIABLES map[string]any) (map[string]any, bool, error) {
	var modified bool

	jsonStr := jsonutils.ToJson(data)

	if jsonStr != "" {
		jsonStr, err := templating.TransformTemplateKeys(jsonStr, '@', func(re *regexp.Regexp, match string) string {
			return re.ReplaceAllStringFunc(match, func(varMatch string) string {
				varName := re.ReplaceAllString(varMatch, "$1")

				return "." + varName
			})
		})

		if err != nil {
			return data, false, err
		}

		normalizedData, err := jsonutils.GetJsonSafe[map[string]any](jsonStr)

		if err == nil {
			data = normalizedData
		}
	}

	templatedData, err := templating.RenderJSON("body", data, VARIABLES)

	if err != nil {
		return data, false, err
	}

	beforeStr := jsonutils.ToJson(data)
	afterStr := jsonutils.ToJson(templatedData)

	log.Dev(beforeStr)
	log.Dev(afterStr)

	modified = beforeStr != afterStr

	return templatedData, modified, nil
}

func TemplatePath(reqUrl *url.URL, VARIABLES any) (string, bool, error) {
	var modified bool

	reqPath, err := url.PathUnescape(reqUrl.Path)

	if err != nil {
		return reqUrl.Path, modified, err
	}

	reqPath, err = templating.RenderNormalizedTemplate("path", reqPath, VARIABLES)

	if err != nil {
		return reqUrl.Path, modified, err
	}

	if reqUrl.Path != reqPath {
		modified = true
	}

	return reqPath, modified, nil
}

func TemplateQuery(reqUrl *url.URL, data map[string]any, VARIABLES any) (string, map[string]any, bool, error) {
	var modified bool

	decodedQuery, _ := url.QueryUnescape(reqUrl.RawQuery)

	templatedQuery, _ := templating.RenderNormalizedTemplate("query", decodedQuery, VARIABLES)

	originalQueryData := reqUrl.Query()

	addedData := query.ParseTypedQuery(templatedQuery, "@")

	for key, val := range addedData {
		data[key] = val

		originalQueryData.Del(key)

		modified = true
	}

	reqRawQuery := originalQueryData.Encode()

	return reqRawQuery, data, modified, nil
}
