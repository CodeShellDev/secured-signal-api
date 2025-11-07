package requestkeys

import "github.com/codeshelldev/secured-signal-api/utils/request"

type Field struct {
	Prefix string
	Key string
}

var BodyPrefix = "@"
var HeaderPrefix = "#"

func Parse(str string) Field {
	prefix := str[:1]
	key := str[1:]

	return Field{
		Prefix: prefix,
		Key: key,
	}
}

func GetByField(field Field, data map[string]any) any {
	key := field.Prefix + field.Key

	return data[key]
}

func PrefixBody(body map[string]any) map[string]any {
	res := map[string]any{}

	for key, value := range body {
		res[BodyPrefix + key] = value
	}

	return res
}

func PrefixHeaders(headers map[string][]string) map[string][]string {
	res := map[string][]string{}

	for key, value := range headers {
		res[HeaderPrefix + key] = value
	}

	return res
}

func PrefixBodyAndHeaders(body map[string]any, headers map[string][]string) (map[string]any, map[string][]string) {
	body = PrefixBody(body)
	headers = PrefixHeaders(headers)

	return body, headers
}

func GetFromBodyAndHeaders(field Field, body map[string]any, headers map[string][]string) any {
	body, headers = PrefixBodyAndHeaders(body, headers)

	switch(field.Prefix) {
	case BodyPrefix:
		return GetByField(field, body)
	case HeaderPrefix:
		return GetByField(field, request.ParseHeaders(headers))
	}

	return nil
}