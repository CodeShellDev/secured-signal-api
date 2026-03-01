package generics

import (
	"errors"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type MatchRule[T any] struct {
	Pattern   T					`koanf:"value"`
	MatchType Enum[MatchType]	`koanf:"matchtype"`
}

type StringMatchRule struct {
	Pattern   string			`koanf:"pattern"`
	MatchType Enum[MatchType]	`koanf:"matchtype"`
}

type MatchType int

const (
	MatchExact MatchType = iota
	MatchEquals
	MatchRegex
	MatchGlob
	MatchContains
	MatchIncludes
	MatchHas
	MatchPrefix
	MatchSuffix
)

func (m MatchType) ParseEnum(str string) (MatchType, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	switch str {
	case "exact":
		return MatchExact, true
	case "equals":
		return MatchEquals, true
	case "regex":
		return MatchRegex, true
	case "glob":
		return MatchGlob, true
	case "contains":
		return MatchContains, true
	case "includes":
		return MatchIncludes, true
	case "has":
		return MatchHas, true
	case "prefix":
		return MatchPrefix, true
	case "suffix":
		return MatchSuffix, true
	default:
		return -1, false
	}
}

func (m MatchType) String() string {
	switch m {
	case MatchExact:
		return "exact"
	case MatchEquals:
		return "equals"
	case MatchRegex:
		return "regex"
	case MatchGlob:
		return "glob"
	case MatchContains:
		return "contains"
	case MatchIncludes:
		return "includes"
	case MatchHas:
		return "has"
	case MatchPrefix:
		return "prefix"
	case MatchSuffix:
		return "suffix"
	default:
		return ""
	}
}

func (r StringMatchRule) Match(str string) (bool, error) {
	rule := MatchRule[string]{
		Pattern: r.Pattern,
		MatchType: r.MatchType,
	}

	return rule.Match(str)
}

func (r StringMatchRule) Test() error {
	rule := MatchRule[string]{
		Pattern: r.Pattern,
		MatchType: r.MatchType,
	}

	return rule.Test()
}

func (r MatchRule[T]) Test() error {
	p := any(r.Pattern)

	switch r.MatchType.Value {
	case MatchEquals:
		_, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}
	case MatchContains:
		_, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}
	case MatchPrefix:
		_, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}
	case MatchSuffix:
		_, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}
	case MatchRegex:
		pStr, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}

		_, err := regexp.Compile(pStr)

		if err != nil {
			return errors.New("could not compile " + pStr + " as regex: " + err.Error())
		}
	case MatchGlob:
		pStr, ok := p.(string)

		if !ok {
			return errors.New("pattern must be string to be able to use match type " + r.MatchType.Value.String())
		}

		_, err := path.Match(pStr, " ")

		if err != nil {
			return errors.New("could not compile " + pStr + " as glob-style pattern: " + err.Error())
		}
	}

	return nil
}

func (r MatchRule[T]) Match(value T) (bool, error) {
	v := any(value)
	p := any(r.Pattern)

	switch r.MatchType.Value {
	case MatchExact:
		return reflect.DeepEqual(v, p), nil
	case MatchEquals:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type equals is only allowed for strings")
		}

		return strings.EqualFold(vStr, pStr), nil
	case MatchContains:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type contains is only allowed for strings")
		}

		return strings.Contains(strings.ToLower(vStr), strings.ToLower(pStr)), nil
	case MatchIncludes:
		vVal := reflect.ValueOf(v)

		if vVal.Kind() == reflect.Slice || vVal.Kind() == reflect.Array {
			pVal := reflect.ValueOf(p)

			for i := 0; i < vVal.Len(); i++ {
				if reflect.DeepEqual(vVal.Index(i).Interface(), pVal.Interface()) {
					return true, nil
				}
			}

			return false, nil
		}

		return false, errors.New("match type includes is not supported for type " + vVal.Kind().String())
	case MatchHas:
		vVal := reflect.ValueOf(v)

		if vVal.Kind() == reflect.Map {
			pVal := reflect.ValueOf(p)

			for _, key := range vVal.MapKeys() {
				if reflect.DeepEqual(key.Interface(), pVal.Interface()) {
					return true, nil
				}
			}

			return false, nil
		}

		return false, errors.New("match type has is only supported for maps")
	case MatchPrefix:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type prefix is only supported for strings")
		}

		return strings.HasPrefix(strings.ToLower(vStr), strings.ToLower(pStr)), nil
	case MatchSuffix:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type suffix is only supported for strings")
		}

		return strings.HasSuffix(strings.ToLower(vStr), strings.ToLower(pStr)), nil
	case MatchRegex:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type regex is only supported for strings")
		}

		re, err := regexp.Compile(pStr)

		if err != nil {
			return false, errors.New("error during regex compilation of " + pStr + ": " + err.Error())
		}
		
		return re.MatchString(vStr), nil
	case MatchGlob:
		vStr, ok1 := v.(string)
		pStr, ok2 := p.(string)

		if !ok1 || !ok2 {
			return false, errors.New("match type glob is only supported for strings")
		}

		match, err := path.Match(pStr, vStr)

		if err != nil {
			return false, errors.New("error during glob-style pattern compilation of " + pStr + ": " + err.Error())
		}

		return match, nil
	default:
		return false, errors.New("unsupported match type")
	}
}