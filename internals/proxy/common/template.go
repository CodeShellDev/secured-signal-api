package common

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/codeshelldev/gotl/pkg/jsonutils"
	queryutils "github.com/codeshelldev/gotl/pkg/query"
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
		res[prefix + key] = value
	}

	return res
}

func cleanHeaders(headers map[string][]string) map[string][]string {
	cleanedHeaders := map[string][]string{}

	for key, value := range headers {
		cleanedKey := strings.ReplaceAll(key, "-", "_")

		cleanedHeaders[cleanedKey] = value
	}

	return cleanedHeaders
}

func GetTemplatedBody(body map[string]any, headers map[string][]string, VARIABLES map[string]any) (map[string]any, bool, error) {
	var modified bool

	bodyCopy := map[string]any{}
	headersCopy := map[string][]string{}

	request.CopyHeaders(headersCopy, headers)
	request.CopyMap(bodyCopy, body)

	headersCopy = cleanHeaders(headersCopy)

	// Normalize `keys.BodyPrefix` + "Var" and `keys.HeaderPrefix` + "Var" to ".header.Var" and ".body.Var"
	normalizedBody, err := normalizeData(requestkeys.BodyPrefix, "body.", bodyCopy)

	if err != nil {
		return bodyCopy, false, err
	}

	normalizedBody, err = normalizeData(requestkeys.HeaderPrefix, "header.", normalizedBody)

	if err != nil {
		return bodyCopy, false, err
	}

	// Prefix Body Data with Body_
	prefixedBody := map[string]any{
		"body": normalizedBody,
	}

	// Prefix Header Data with Header_
	prefixedHeaders := map[string]any{
		"header": request.ParseHeaders(headersCopy),
	}

	variables := map[string]any{}

	request.CopyMap(variables, VARIABLES)
	request.CopyMap(variables, prefixedBody)
	request.CopyMap(variables, prefixedHeaders)

	templatedData, err := templating.RenderJSON(normalizedBody, variables)

	if err != nil {
		return bodyCopy, false, err
	}

	beforeStr := jsonutils.ToJson(bodyCopy)
	afterStr := jsonutils.ToJson(templatedData)

	modified = beforeStr != afterStr

	return templatedData, modified, nil
}

func TemplatePath(path string, VARIABLES any) (string, error) {
	reqPath, err := url.PathUnescape(path)

	if err != nil {
		return path, err
	}

	reqPath, err = templating.RenderNormalizedTemplate("path", reqPath, VARIABLES)

	if err != nil {
		return path, err
	}

	return reqPath, nil
}

func InjectPathIntoBody(path string, data map[string]any) (string, bool) {
	var modified bool

	parts := strings.Split(path, "/")
	newParts := []string{}

	for _, part := range parts {
		newParts = append(newParts, part)
		
		keyValuePair := strings.SplitN(part, "=", 2)

		if len(keyValuePair) != 2 {
			continue
		}

		keyWithoutPrefix, match := strings.CutPrefix(keyValuePair[0], requestkeys.BodyPrefix)
		
		if !match {
			continue
		}

		value := stringutils.ToType(keyValuePair[1])

		data[keyWithoutPrefix] = value
		modified = true

		newParts = newParts[:len(newParts) - 1]
	}

	return strings.Join(newParts, "/"), modified
}

func TemplateQuery(rawQuery string, VARIABLES any) (string, error) {
	decodedQuery, _ := url.QueryUnescape(rawQuery)

	templatedQuery, err := templating.RenderNormalizedTemplate("query", decodedQuery, VARIABLES)

	return templatedQuery, err
}

func InjectQueryIntoBody(query url.Values, data map[string]any) bool {
	var modified bool

	decodedQuery, _ := url.QueryUnescape(query.Encode())

	parsedQuery, _ := queryutils.ParseTypedQuery(decodedQuery)

	for key, val := range parsedQuery {
		keyWithoutPrefix, match := strings.CutPrefix(key, requestkeys.BodyPrefix)

		if !match {
			continue
		}

		data[keyWithoutPrefix] = val

		query.Del(key)

		modified = true
	}

	return modified
}