package structure

import (
	"errors"
	"strconv"

	"github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
)

func (e Endpoints) Test(path string) error {
	if len(e.Allowed) != 0 {
		for i, rule := range e.Allowed {
			err := rule.Test()

			if err != nil {
				return errors.New(strconv.Itoa(i) + ": " + err.Error())
			}
		}
	}

	if len(e.Blocked) != 0 {
		for i, rule := range e.Blocked {
			err := rule.Test()

			if err != nil {
				return errors.New(strconv.Itoa(i) + ": " + err.Error())
			}
		}
	}

	return nil
}

func (p FieldPolicies) Test(path string) error {
	policies := map[string]custom.FPolicies(p.Value.Compile())

	for field, policy := range policies {
		for _, item := range policy.Allowed {
			err := item.MatchRule.Test()

			if err != nil {
				return errors.New("'" + field + "': " + err.Error())
			}
		}

		for _, item := range policy.Blocked {
			err := item.MatchRule.Test()

			if err != nil {
				return errors.New("'" + field + "': " + err.Error())
			}
		}
	}

	return nil
}