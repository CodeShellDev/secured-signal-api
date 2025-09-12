package safestrings

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils"
)

func ToType(str string) any {
	cleaned := strings.TrimSpace(str)

    //* Try JSON
	if IsEnclosedBy(cleaned, `[`, `]`) || IsEnclosedBy(cleaned, `{`, `}`) {
		data, err := utils.GetJsonSafe[any](str)

		if data != nil && err == nil {
			return data
		}
	}

	//* Try String Slice
	if Contains(str, ",") && IsEnclosedBy(cleaned, `[`, `]`) {
		bracketsless := strings.ReplaceAll(str, "[", "")
		bracketsless = strings.ReplaceAll(bracketsless, "]", "")

		data := ToArray(bracketsless)

		if data != nil {
			if len(data) > 0 {
				return data
			}
		}
	}

	//* Try Number
	if !strings.HasPrefix(cleaned, "+") {
		intValue, intErr := strconv.Atoi(cleaned)

		if intErr == nil {
			return intValue
		}
	}

    return str
}

func Contains(str string, match string) bool {
    return !IsEscaped(str, match)
}

// Checks if a string is Enclosed by `char` and are not Escaped
func IsEnclosedBy(str string, charA, charB string) bool {
	if NeedsEscapeForRegex(rune(charA[0])) {
		charA = `\` + charA
	}

	if NeedsEscapeForRegex(rune(charB[0])) {
		charB = `\` + charB
	}

	regexStr := `(^|[^\\])(\\\\)*(` + charA + `)(.*?)(^|[^\\])(\\\\)*(` + charB + ")"

 	re := regexp.MustCompile(regexStr)

	matches := re.FindAllStringSubmatchIndex(str, -1)

	filtered := [][]int{}

	for _, match := range matches {
		start := match[len(match)-2]
		end := match[len(match)-1]
		char := str[start:end]

		if char != `\` {
			filtered = append(filtered, match)
		}
	}

	return len(filtered) > 0
}

// Checks if a string is completly Escaped with `\`
func IsEscaped(str string, char string) bool {
	if NeedsEscapeForRegex(rune(char[0])) {
		char = `\` + char
	}

	regexStr := `(^|[^\\])(\\\\)*(` + char + ")"

	re := regexp.MustCompile(regexStr)

	matches := re.FindAllStringSubmatchIndex(str, -1)

	filtered := [][]int{}

	for _, match := range matches {
		start := match[len(match)-2]
		end := match[len(match)-1]
		char := str[start:end]

		if char != `\` {
			filtered = append(filtered, match)
		}
	}

	return len(filtered) == 0
}

func NeedsEscapeForRegex(char rune) bool {
	special := `.+*?()|[]{}^$\\`

	return strings.ContainsRune(special, char)
}

func ToArray(sliceStr string) []string {
    if sliceStr == "" {
        return nil
    }

    rawItems := strings.Split(sliceStr, ",")
    items := make([]string, 0, len(rawItems))

    for _, item := range rawItems {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" {
            items = append(items, trimmed)
        }
    }

    return items
}