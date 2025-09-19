package templating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/codeshelldev/secured-signal-api/utils/logger"
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

func RenderJSON(name string, data map[string]any, variables any) (map[string]any, error) {
	data, err := RenderJSONTemplate(name + ":json_path", data, data)

	if err != nil {
		return data, err
	}

	data, err = RenderJSONTemplate(name + ":variables", data, variables)

	if err != nil {
		return data, err
	}

	return data, nil
}

func RenderJSONTemplate(name string, data map[string]any, variables any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	tmplStr := string(jsonBytes)

	re, err := regexp.Compile(`{{\s*\.([a-zA-Z0-9_.]+)\s*}}`)

	// Add normalize() to be able to remove Quotes from Arrays
	if err == nil {
    	tmplStr = re.ReplaceAllString(tmplStr, "{{normalize .$1}}")
	}

	templt := CreateTemplateWithFunc(name, template.FuncMap{
        "normalize": normalizeJSON,
    })

	jsonStr, err := ParseTemplate(templt, tmplStr, variables)

	if err != nil {
		return nil, err
	}

	logger.Dev("after template:\n" + jsonStr)

	jsonStr = cleanQuotedPairsJSON(jsonStr)

	// Remove the Quotes around "<<[item1,item2]>>"
	re, err = regexp.Compile(`"<<(.*?)>>"`)

	if err != nil {
		return nil, err
	}

	jsonStr = re.ReplaceAllString(jsonStr, "$1")

	err = json.Unmarshal([]byte(jsonStr), &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func RenderNormalizedTemplate(name string, tmplStr string, variables any) (string, error) {
	re, err := regexp.Compile(`{{\s*\.(\w+)\s*}}`)

	// Add normalize() to normalize arrays to [item1,item2]
	if err == nil {
    	tmplStr = re.ReplaceAllString(tmplStr, "{{normalize .$1}}")
	}

	templt := CreateTemplateWithFunc(name, template.FuncMap{
        "normalize": normalize,
    })

	return ParseTemplate(templt, tmplStr, variables)
}