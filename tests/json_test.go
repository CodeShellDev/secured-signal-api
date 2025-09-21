package tests

import (
	"testing"

	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	templating "github.com/codeshelldev/secured-signal-api/utils/templating"
)

func TestJsonTemplating(t *testing.T) {
	variables := map[string]any{
		"array": []string{
			"item0",
			"item1",
		},
		"key": "val",
		"int": 4,
	}

	json := `
	{
		"multiple": "{{.key}}, {{.int}}",
		"dict": { "key": "{{.key}}" },
		"dictArray": [
			{ "key": "{{.key}}" },
			{ "key": "{{.array}}" }
		],
		"key1": "{{.array}}",
		"key2": "{{.int}}"
	}`

	data := jsonutils.GetJson[map[string]any](json)

	expected := map[string]any{
		"multiple": "val, 4",
		"dict": map[string]any{
			"key": "val",
		},
		"dictArray": []any{
			map[string]any{"key": "val"},
			map[string]any{"key": []any{ "item0", "item1" }},
		},
		"key1": []any{ "item0", "item1" },
		"key2": 4,
	}

	got, err := templating.RenderDataKeyTemplateRecursive("", data, variables)

	if err != nil {
		t.Error("Error Templating JSON:\n", err.Error())
	}

	expectedStr := jsonutils.ToJson(expected)
	gotStr := jsonutils.ToJson(got)

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

	data := jsonutils.GetJson[map[string]any](json)

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

		got, ok := jsonutils.GetByPath(key, data)

		if !ok || got.(string) != expected {
			t.Error("Expected: ", key, " == ", expected, "; Got: ", got)
		}
	}
}