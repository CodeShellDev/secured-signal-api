package templating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"strings"
	"text/template"

	"github.com/codeshelldev/secured-signal-api/utils/stringutils"
)

func normalize(value any) string {
	switch str := value.(type) {
		case []string:
			return "[" + strings.Join(str, ",") + "]"
		case []any:
			items := make([]string, len(str))

			for i, item := range str {
				items[i] = fmt.Sprintf("%v", item)
			}

			return "[" + strings.Join(items, ",") + "]"
		default:
			return fmt.Sprintf("%v", value)
	}
}

func normalizeJSON(value any) string {
	if value == nil {
		return ""
	}

	switch value.(type) {
		case []any, []string, map[string]any, int, float64, bool:
			object, _ := json.Marshal(value)

			if string(object) == "{}" {
				return value.(string)
			}

			return "<<" + string(object) + ">>"

		default:
			return value.(string)
    }
}

func cleanQuotedPairsJSON(s string) string {
	quoteRe, err := regexp.Compile(`"([^"]*?)"`)

	if err != nil {
		return s
	}

	pairRe, err := regexp.Compile(`<<([^<>]+)>>`)

	if err != nil {
		return s
	}

	return quoteRe.ReplaceAllStringFunc(s, func(container string) string {
		inner := container[1 : len(container)-1] // remove quotes

		matches := pairRe.FindAllStringSubmatchIndex(inner, -1)

		// ONE pair which fills whole ""
		if len(matches) == 1 && matches[0][0] == 0 && matches[0][1] == len(inner) {
			return container // keep <<...>> untouched
		}

		// MULTIPLE pairs || that do not fill whole ""
		inner = pairRe.ReplaceAllString(inner, "$1")
		inner = strings.ReplaceAll(inner, `"`, `'`)

		return `"` + inner + `"`
	})
}

func AddTemplateFunc(tmplStr string, funcName string) (string, error) {
	return TransformTemplateKeys(tmplStr, '.', func(re *regexp.Regexp, match string) string {
		reSimple, _ := regexp.Compile(`{{\s*\.[a-zA-Z0-9_.]+\s*}}`)

		if !reSimple.MatchString(match) {
			return match
		}

		return re.ReplaceAllStringFunc(match, func(varMatch string) string {
			varName := re.ReplaceAllString(varMatch, ".$1")

			return strings.ReplaceAll(varMatch, varName, "(" + funcName + " " + varName + ")")
		})
	})
}

func TransformTemplateKeys(tmplStr string, tmplKey rune, transform func(varRegex *regexp.Regexp, m string) string) (string, error) {
	re, err := regexp.Compile(`{{[^}]+}}`)

	if err != nil {
		return tmplStr, err
	}

	varRe, err := regexp.Compile(`\` + string(tmplKey) + `([a-zA-Z0-9_.]+)`)

	if err != nil {
		return tmplStr, err
	}

	tmplStr = re.ReplaceAllStringFunc(tmplStr, func(match string) string {
		return transform(varRe, match)
	})

	return tmplStr, nil
}

func ParseTemplate(templt *template.Template, tmplStr string, variables any) (string, error) {
	tmpl, err := templt.Parse(tmplStr)

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, variables)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderTemplate(name string, tmplStr string, variables any) (string, error) {
	templt := template.New(name)

	return ParseTemplate(templt, tmplStr, variables)
}

func CreateTemplateWithFunc(name string, funcMap template.FuncMap) (*template.Template) {
	return template.New(name).Funcs(funcMap)
}

func RenderJSON(name string, data map[string]any, variables map[string]any) (map[string]any, error) {
	combinedData := data

	maps.Copy(combinedData, variables)

	data, err := RenderJSONTemplate(name, data, combinedData)

	if err != nil {
		return data, err
	}

	return data, nil
}

func RenderDataKeyTemplateRecursive(key any, value any, variables map[string]any) (any, error) {
	var err error

	strKey, isStr := key.(string)

	if !isStr {
		strKey = "!string"
	}

	switch typedValue := value.(type) {
		case map[string]any:
			data := map[string]any{}

			for mapKey, mapValue := range typedValue {
				var templatedValue any

				templatedValue, err = RenderDataKeyTemplateRecursive(mapKey, mapValue, variables)

				if err != nil {
					return mapValue, err
				}

				data[mapKey] = templatedValue
			}

			return data, err

		case []any:
			data := []any{}

			for arrayIndex, arrayValue := range typedValue {
				var templatedValue any

				templatedValue, err = RenderDataKeyTemplateRecursive(arrayIndex, arrayValue, variables)

				if err != nil {
					return arrayValue, err
				}

				data = append(data, templatedValue)
			}

			return data, err
		
		case string:
			templt := CreateTemplateWithFunc("json:" + strKey, template.FuncMap{
				"normalize": normalize,
			})

			tmplStr, _ := AddTemplateFunc(typedValue, "normalize")

			templatedValue, err := ParseTemplate(templt, tmplStr, variables)

			if err != nil {
				return typedValue, err
			}

			templateRe, err := regexp.Compile(`{{[^}]+}}`)

			if err == nil {
				nonWhitespaceRe, err := regexp.Compile(`(\S+)`)

				if err == nil {
					filtered := nonWhitespaceRe.ReplaceAllString(tmplStr, "")

					if !templateRe.MatchString(filtered) {
						return stringutils.ToType(templatedValue), err
					}
				}
			}
		
			return templatedValue, err
		
		default:
			return typedValue, err
	}
}

func RenderJSONTemplate(name string, data map[string]any, variables map[string]any) (map[string]any, error) {
	res, err := RenderDataKeyTemplateRecursive("", data, variables)

	mapRes, ok := res.(map[string]any)

	if !ok {
		return data, err
	}

	return mapRes, err
}

func RenderNormalizedTemplate(name string, tmplStr string, variables any) (string, error) {
	tmplStr, err := AddTemplateFunc(tmplStr, "normalize")

	if err != nil {
		return tmplStr, err
	}

	templt := CreateTemplateWithFunc(name, template.FuncMap{
        "normalize": normalize,
    })

	return ParseTemplate(templt, tmplStr, variables)
}
