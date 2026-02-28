package structure

import (
	g "github.com/codeshelldev/secured-signal-api/internals/config/structure/generics"
)
type StringMatchList []g.StringMatchRule

func (m StringMatchList) FindMatchRule(str string) (g.StringMatchRule, error) {
	for _, rule := range m {
		ok, err := rule.Match(str)

		if ok {
			return rule, err
		}
	}

	return g.StringMatchRule{}, nil
}

func (m StringMatchList) Match(str string) (bool, error) {
	rule, err := m.FindMatchRule(str)

	if err != nil {
		return false, err
	}

	return rule.Pattern != "", nil
}