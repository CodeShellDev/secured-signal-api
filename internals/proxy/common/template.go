package common

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/codeshelldev/gotl/pkg/jsonutils"
	"github.com/codeshelldev/gotl/pkg/query"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/gotl/pkg/stringutils"
	"github.com/codeshelldev/gotl/pkg/templating"
	"github.com/codeshelldev/secured-signal-api/utils/requestkeys"
)

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

	request.CopyMap(variables, VARIABLES)
	request.CopyMap(variables, prefixedBody)
	request.CopyMap(variables, prefixedHeaders)

	templatedData, err := templating.RenderJSON(normalizedBody, variables)

	if err != nil {
		return body, false, err
	}

	beforeStr := jsonutils.ToJson(body)
	afterStr := jsonutils.ToJson(templatedData)

	modified = beforeStr != afterStr

	return templatedData, modified, nil
}

func TemplatePath(reqUrl *url.URL, data map[string]any, VARIABLES any) (string, map[string]any, bool, bool, error) {
	var modified bool
	var modifiedBody bool

	reqPath, err := url.PathUnescape(reqUrl.Path)

	if err != nil {
		return reqUrl.Path, data, false, false, err
	}

	reqPath, err = templating.RenderNormalizedTemplate("path", reqPath, VARIABLES)

	if err != nil {
		return reqUrl.Path, data, false, false, err
	}

	parts := strings.Split(reqPath, "/")
	newParts := []string{}

	for _, part := range parts {
		newParts = append(newParts, part)
		
		keyValuePair := strings.SplitN(part, "=", 2)

		if len(keyValuePair) != 2 {
			continue
		}

		keyWithoutPrefix, match := strings.CutPrefix(keyValuePair[0], "@")
		
		if !match {
			continue
		}

		value := stringutils.ToType(keyValuePair[1])

		data[keyWithoutPrefix] = value
		modifiedBody = true

		newParts = newParts[:len(newParts) - 1]
	}

	reqPath = strings.Join(newParts, "/")

	if reqUrl.Path != reqPath {
		modified = true
	}

	return reqPath, data, modified, modifiedBody, nil
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
