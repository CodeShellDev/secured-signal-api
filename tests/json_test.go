package tests

import (
	"testing"

	"github.com/codeshelldev/secured-signal-api/utils"
	"github.com/codeshelldev/secured-signal-api/utils/templating"
)

func TestJsonTemplating(t *testing.T) {
	variables := map[string]interface{}{
		"array": []string{
			"item0",
			"item1",
		},
		"key": "val",
		"int": 4,
	}

	json := `
	{
		"dict": { "key": "{{.key}}" },
		"dictArray": [
			{ "key": "{{.key}}" },
			{ "key": "{{.array}}" }
		],
		"key1": "{{.array}}",
		"key2": "{{.int}}"
	}`

	data := utils.GetJson[map[string]interface{}](json)

	expected := map[string]interface{}{
		"dict": map[string]interface{}{
			"key": "val",
		},
		"dictArray": []interface{}{
			map[string]interface{}{"key": "val"},
			map[string]interface{}{"key": []interface{}{ "item0", "item1" }},
		},
		"key1": []interface{}{ "item0", "item1" },
		"key2": 4,
	}

	got, err := templating.RenderJSONTemplate("json", data, variables)

	if err != nil {
		t.Error("Error Templating JSON: ", err.Error())
	}

	expectedStr := utils.ToJson(expected)
	gotStr := utils.ToJson(got)

	if expectedStr != gotStr {
		t.Error("\nExpected: ", expectedStr, "\nGot: ", gotStr)
	}
}

func TestJsonPath(t *testing.T) {
	json := `
	{
		"dict": { "key": "value" },
		"dictArray": [
			{ "key": "value0" },
			{ "key": "value1" }
		],
		"array": [
			"item0",
			"item1"
		],
		"key": "val"
	}`

	data := utils.GetJson[map[string]interface{}](json)

	cases := []struct{
		key 	 string
		expected string
	}{
		{
			key: "key",
			expected: "val",
		},
		{
			key: "dict.key",
			expected: "value",
		},
		{
			key: "dictArray[0].key",
			expected: "value0",
		},
		{
			key: "dictArray[1].key",
			expected: "value1",
		},
		{
			key: "array[0]",
			expected: "item0",
		},
		{
			key: "array[1]",
			expected: "item1",
		},
	}

	for _, c := range cases {
		key := c.key
		expected := c.expected

		got, ok := utils.GetByPath(key, data)

		if !ok || got.(string) != expected {
			t.Error("Expected: ", key, " == ", expected, "; Got: ", got)
		}
	}
}