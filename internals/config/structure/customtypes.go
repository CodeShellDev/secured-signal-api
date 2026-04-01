package structure

import (
	t "github.com/codeshelldev/gotl/pkg/configutils/types"
	c "github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
	g "github.com/codeshelldev/secured-signal-api/internals/config/structure/generics"
	"github.com/go-viper/mapstructure/v2"
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

type FieldPolicies struct {
    Value *t.Comp[c.RFieldPolicies, c.FieldPolicies]
}

func (f *FieldPolicies) UnmarshalMapstructure(raw any) error {
    var comp *t.Comp[c.RFieldPolicies, c.FieldPolicies]

    err := mapstructure.Decode(raw, &comp)

	if err != nil {
		return err
	}

	if comp == nil {
		f.Value = comp
		return nil
	}

	comp.Compile()

	f.Value = comp

	return nil
}