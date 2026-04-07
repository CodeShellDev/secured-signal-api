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

func transformPrefixInto(tmplStr string, prefix string, to string) string {	
	re, err := regexp.Compile(`{{([^{}]+)}}`)

	if err != nil {
		return tmplStr
	}

	varRe, err := regexp.Compile(string(prefix) + `("?[a-zA-Z0-9_.]+"?)`)

	if err != nil {
		return tmplStr
	}

	transformed := re.ReplaceAllStringFunc(tmplStr, func(match string) string {
		return varRe.ReplaceAllStringFunc(match, func(varMatch string) string {
			varName := varRe.ReplaceAllString(varMatch, "$1")

			return "." + to + "." + varName
		})
	})

	return transformed
}

func normalizeData(fromPrefix, toPrefix string, data map[string]any) (map[string]any, error) {
	jsonStr := jsonutils.ToJson(data)

	if jsonStr != "" {
		normalizedTemplate := transformPrefixInto(jsonStr, fromPrefix, toPrefix)

		normalizedData, err := jsonutils.GetJsonSafe[map[string]any](normalizedTemplate)

		if err == nil {
			data = normalizedData
		}
	}

	return data, nil
}

func normalizeHeaders(headers map[string][]string) map[string][]string {
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

	headersCopy = normalizeHeaders(headersCopy)

	// Normalize `keys.BodyPrefix` + "Var" and `keys.HeaderPrefix` + "Var" to ".headers.Var" and ".body.Var"
	normalizedBody, err := normalizeData(requestkeys.BodyPrefix, "body", bodyCopy)

	if err != nil {
		return bodyCopy, false, err
	}

	normalizedBody, err = normalizeData(requestkeys.HeaderPrefix, "headers", normalizedBody)

	if err != nil {
		return bodyCopy, false, err
	}

	// Prefix Body Data with Body_
	nestedBody := map[string]any{
		"body": normalizedBody,
	}

	// Prefix Header Data with Header_
	nestedHeaders := map[string]any{
		"headers": request.ParseHeaders(headersCopy),
	}

	variables := map[string]any{}

	request.CopyMap(variables, VARIABLES)
	request.CopyMap(variables, nestedBody)
	request.CopyMap(variables, nestedHeaders)

	templatedData, err := templating.TemplateData(normalizedBody, variables)

	if err != nil {
		return bodyCopy, false, err
	}

	beforeStr := jsonutils.ToJson(bodyCopy)
	afterStr := jsonutils.ToJson(templatedData)

	modified = beforeStr != afterStr

	return templatedData.(map[string]any), modified, nil
}

func TemplatePath(path string, VARIABLES map[string]any) (string, error) {
	reqPath, err := url.PathUnescape(path)

	if err != nil {
		return path, err
	}

	templt, err := templating.CreateNormalizedTemplateFromString("path", reqPath)

	if err != nil {
		return path, err
	}

	templated, err := templating.ExecuteTemplate(templt, VARIABLES)

	if err != nil {
		return path, err
	}

	return templated, nil
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

func TemplateQuery(rawQuery string, VARIABLES map[string]any) (string, error) {
	query, _ := url.ParseQuery(rawQuery)

	newQuery := url.Values{}

	for key, value := range query {
		templt, err := templating.CreateNormalizedTemplateFromString("query", value[0])

		if err != nil {
			return rawQuery, err
		}

		templated, err := templating.ExecuteTemplate(templt, VARIABLES)

		if err != nil {
			return rawQuery, err
		}

		newQuery.Set(key, templated)
	}

	return newQuery.Encode(), nil
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