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
    switch value.(type) {
		case []any, []string, map[string]any, int, float64, bool:
			object, _ := json.Marshal(value)

			return "<<" + string(object) + ">>"

		default:
			return value.(string)
    }
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

func RenderJSONTemplate(name string, data map[string]any, variables any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		logger.Dev("72"+err.Error())
		return nil, err
	}

	tmplStr := string(jsonBytes)

	re, err := regexp.Compile(`{{\s*\.(\w+)\s*}}`)

	// Add normalize() to be able to remove Quotes from Arrays
	if err == nil {
    	tmplStr = re.ReplaceAllString(tmplStr, "{{normalize .$1}}")
	}

	templt := CreateTemplateWithFunc(name, template.FuncMap{
        "normalize": normalizeJSON,
    })

	jsonStr, err := ParseTemplate(templt, tmplStr, variables)

	if err != nil {
		logger.Dev("92:"+err.Error())
		return nil, err
	}

	// Remove the Quotes around "<<[item1,item2]>>"
	re, err = regexp.Compile(`"<<(.*?)>>"`)

	if err != nil {
		logger.Dev("100:"+err.Error())
		return nil, err
	}

	jsonStr = re.ReplaceAllString(jsonStr, "$1")

	err = json.Unmarshal([]byte(jsonStr), &data)

	if err != nil {
		logger.Dev("109:"+err.Error())
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