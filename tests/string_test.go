package tests

import (
	"reflect"
	"testing"

	"github.com/codeshelldev/secured-signal-api/utils/safestrings"
)

func TestStringEscaping(t *testing.T) {
	str1 := `\#`

	res1 := safestrings.IsEscaped(str1, "#")

	if !res1 {
		t.Error("Expected: ", str1, " == true", "; Got: ", str1, " == ", res1)
	}

	str2 := "#"

	res2 := safestrings.IsEscaped(str2, "#")

	if res2 {
		t.Error("Expected: ", str2, " == false", "; Got: ", str2, " == ", res2)
	}

	str3 := `#\#`

	res3 := safestrings.Contains(str3, "#")

	if !res3 {
		t.Error("Expected: ", str3, " == true", "; Got: ", str3, " == ", res3)
	}
}

func TestStringEnclosement(t *testing.T) {
	str1 := "[enclosed]"

	res1 := safestrings.IsEnclosedBy(str1, `[`, `]`)

	if !res1 {
		t.Error("Expected: ", str1, " == true", "; Got: ", str1, " == ", res1)
	}

	str2 := `\[enclosed]`

	res2 := safestrings.IsEnclosedBy(str2, `[`, `]`)

	if res2 {
		t.Error("Expected: ", str2, " == false", "; Got: ", str2, " == ", res2)
	}
}

func TestStringToType(t *testing.T) {
	str1 := `[item1,item2]`

	res1 := safestrings.ToType(str1)

	if reflect.TypeOf(res1) != reflect.TypeFor[[]string]() {
		t.Error("Expected: ", str1, " == []string", "; Got: ", str1, " == ", reflect.TypeOf(res1))
	}

	str2 := `1`

	res2 := safestrings.ToType(str2)

	if reflect.TypeOf(res2) != reflect.TypeFor[int]() {
		t.Error("Expected: ", str2, " == int", "; Got: ", str2, " == ", reflect.TypeOf(res2))
	}

	str3 := `{ "key": "value" }`

	res3 := safestrings.ToType(str3)

	if reflect.TypeOf(res3) != reflect.TypeFor[map[string]any]() {
		t.Error("Expected: ", str3, " == map[string]any", "; Got: ", str3, " == ", reflect.TypeOf(res3))
	}
}