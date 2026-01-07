package middlewares

import (
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	jsonutils "github.com/codeshelldev/gotl/pkg/jsonutils"
	query "github.com/codeshelldev/gotl/pkg/query"
	request "github.com/codeshelldev/gotl/pkg/request"
	templating "github.com/codeshelldev/gotl/pkg/templating"
	"github.com/codeshelldev/secured-signal-api/utils/requestkeys"
)

var Template Middleware = Middleware{
	Name: "Template",
	Use: templateHandler,
}

func templateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)
		
		variables := conf.SETTINGS.MESSAGE.VARIABLES

		if variables == nil {
			variables = getConfig("").SETTINGS.MESSAGE.VARIABLES
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty {
			var modified bool

			headerData := request.GetReqHeaders(req)

			bodyData, modified, err = TemplateBody(body.Data, headerData, variables)

			if err != nil {
				logger.Error("Error Templating JSON: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.RawQuery != "" {
			var modified bool

			req.URL.RawQuery, bodyData, modified, err = TemplateQuery(req.URL, bodyData, variables)

			if err != nil {
				logger.Error("Error Templating Query: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if modifiedBody {
			body.Data = bodyData

			err := body.Write(req)

			if err != nil {
				logger.Error("Could not write to Request Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			logger.Debug("Applied Body Templating: ", body.Data)
		}

		if req.URL.Path != "" {
			var modified bool

			req.URL.Path, modified, err = TemplatePath(req.URL, variables)

			if err != nil {
				logger.Error("Error Templating Path: ", err.Error())
			}

			if modified {
				logger.Debug("Applied Path Templating: ", req.URL.Path)
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

func prefixData(prefix string, data map[string]any) map[string]any {
	res := map[string]any{}

	for key, value := range data {
		res[prefix+key] = value
	}

	return res
}

func cleanHeaders(headers map[string][]string) map[string][]string {
	cleanedHeaders := map[string][]string{}

	for key, value := range headers {
		cleanedKey := strings.ReplaceAll(key, "-", "_")

		cleanedHeaders[cleanedKey] = value
	}

	authHeader, ok := cleanedHeaders["Authorization"]

	if !ok {
		authHeader = []string{"UNKNOWN REDACTED"}
	}

	cleanedHeaders["Authorization"] = []string{strings.Split(authHeader[0], ` `)[0] + " REDACTED"}

	return cleanedHeaders
}

func TemplateBody(body map[string]any, headers map[string][]string, VARIABLES map[string]any) (map[string]any, bool, error) {
	var modified bool

	headers = cleanHeaders(headers)

	// Normalize `keys.BodyPrefix` + "Var" and `keys.HeaderPrefix` + "Var" to "".header_key_Var" and ".body_key_Var"
	normalizedBody, err := normalizeData(requestkeys.BodyPrefix, "body_key_", body)

	if err != nil {
		return body, false, err
	}

	normalizedBody, err = normalizeData(requestkeys.HeaderPrefix, "header_key_", normalizedBody)

	if err != nil {
		return body, false, err
	}

	// Prefix Body Data with body_key_
	prefixedBody := prefixData("body_key_", normalizedBody)

	// Prefix Header Data with header_key_
	prefixedHeaders := prefixData("header_key_", request.ParseHeaders(headers))

	variables := map[string]any{}

	maps.Copy(variables, VARIABLES)

	maps.Copy(variables, prefixedBody)
	maps.Copy(variables, prefixedHeaders)

	templatedData, err := templating.RenderJSON(normalizedBody, variables)

	if err != nil {
		return body, false, err
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

	addedData, _ := query.ParseTypedQuery(templatedQuery)

	for key, val := range addedData {
		keyWithoutPrefix, match := strings.CutPrefix(key, "@")

		if !match {
			continue
		}

		data[keyWithoutPrefix] = val

		originalQueryData.Del(key)

		modified = true
	}

	reqRawQuery := originalQueryData.Encode()

	return reqRawQuery, data, modified, nil
}
