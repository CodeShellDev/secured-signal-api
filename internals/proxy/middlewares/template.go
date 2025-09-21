package middlewares

import (
	"bytes"
	"io"
	"maps"
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

			headerData := request.GetReqHeaders(req)

			bodyData, modified, err = TemplateBody(body.Data, headerData, variables)

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

func normalizeData(fromPrefix, toPrefix string, data map[string]any) (map[string]any, error) {
	jsonStr := jsonutils.ToJson(data)

	if jsonStr != "" {
		toVar, err := templating.TransformTemplateKeys(jsonStr, fromPrefix, func(re *regexp.Regexp, match string) string {
			return re.ReplaceAllStringFunc(match, func(varMatch string) string {
				varName := re.ReplaceAllString(varMatch, "$1")

				return "." + toPrefix + varName
			})
		})

		if err != nil {
			return data, err
		}

		jsonStr = toVar

		normalizedData, err := jsonutils.GetJsonSafe[map[string]any](jsonStr)

		if err == nil {
			data = normalizedData
		}
	}

	return data, nil
}

func prefixData(prefix string, data map[string]any) (map[string]any) {
	res := map[string]any{}

	for key, value := range data {
		res[prefix + key] = value
	}

	return res
}

func TemplateBody(body map[string]any, headers map[string]any, VARIABLES map[string]any) (map[string]any, bool, error) {
	var modified bool

	// Normalize #Var and @Var to .header_key_Var and .body_key_Var
	bodyData, err := normalizeData("@", "body_key_", body)

	log.Dev("Normalized:\n", jsonutils.ToJson(bodyData))

	if err != nil {
		return body, false, err
	}

	headerData, err := normalizeData("#", "header_key_", headers)

	if err != nil {
		return body, false, err
	}

	// Prefix Body Data with body_key_
	bodyData = prefixData("body_key_", bodyData)

	// Prefix Header Data with header_key_
	headerData = prefixData("header_key_", headerData)

	variables := VARIABLES
	
	maps.Copy(variables, bodyData)
	maps.Copy(variables, headerData)

	log.Dev("Body:\n", jsonutils.ToJson(bodyData))
	log.Dev("Headers:\n", jsonutils.ToJson(headerData))

	templatedData, err := templating.RenderJSON("body", body, variables)

	if err != nil {
		return bodyData, false, err
	}

	beforeStr := jsonutils.ToJson(body)
	afterStr := jsonutils.ToJson(templatedData)

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
