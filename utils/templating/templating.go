package templating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

func normalize(value interface{}) string {
	switch str := value.(type) {
		case []string:
			return "[" + strings.Join(str, ",") + "]"

		case []interface{}:
			items := make([]string, len(str))

			for i, item := range str {
				items[i] = fmt.Sprintf("%v", item)
			}

			return "[" + strings.Join(items, ",") + "]"
		default:
			return fmt.Sprintf("%v", value)
	}
}

func normalizeJSON(value interface{}) string {
	jsonBytes, err := json.Marshal(value)

	if err != nil {
		return "INVALID:JSON"
	}

	return "<<" + string(jsonBytes) + ">>"
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

func RenderJSONTemplate(name string, data map[string]interface{}, variables any) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
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
		return nil, err
	}

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