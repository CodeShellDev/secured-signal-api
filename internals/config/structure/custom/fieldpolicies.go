package custom

import (
	"errors"
	"strings"

	g "github.com/codeshelldev/secured-signal-api/internals/config/structure/generics"
	"github.com/go-viper/mapstructure/v2"
)

type FPolicyAction int

const (
	FPolicyActionBlock = iota
	FPolicyActionAllow
)

func (m FPolicyAction) ParseEnum(str string) (FPolicyAction, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	switch str {
	case "block":
		return FPolicyActionBlock, true
	case "allow":
		return FPolicyActionAllow, true
	default:
		return -1, false
	}
}

func (m FPolicyAction) String() string {
	switch m {
	case FPolicyActionBlock:
		return "block"
	case FPolicyActionAllow:
		return "allow"
	default:
		return ""
	}
}

type RFPolicy struct {
	Action		g.Enum[FPolicyAction]	`koanf:"action"`
	Value		any						`koanf:"value"`
	MatchType 	g.Enum[g.MatchType] 	`koanf:"matchtype"`
}

type RFieldPolicies map[string][]RFPolicy

func (r *RFieldPolicies) UnmarshalMapstructure(raw any) error {
	rawMap, ok := raw.(map[string]any)

	if !ok {
		return errors.New("expected map input")
	}

	result := make(RFieldPolicies, len(rawMap))

	for key, val := range rawMap {
		var policies []RFPolicy

		err := mapstructure.Decode(val, &policies)

		if err != nil {
			return err
		}

		result[key] = policies
	}

	*r = result

	return nil
}


type FPolicy struct {
	Action    FPolicyAction
	MatchRule g.MatchRule[any]
}

type FPolicies struct {
	Allowed []FPolicy
	Blocked []FPolicy
}

type FieldPolicies map[string]FPolicies

func (r RFPolicy) Compile() FPolicy {
	return FPolicy{
		Action: r.Action.Value,
		MatchRule: g.MatchRule[any]{
			MatchType: r.MatchType,
			Pattern:   r.Value,
		},
	}
}

func (r RFieldPolicies) Compile() FieldPolicies {
	out := make(FieldPolicies)

	for field, policies := range r {
		var allowed []FPolicy
		var blocked []FPolicy

		for _, p := range policies {
			fp := p.Compile()
			
			if fp.Action == FPolicyActionAllow {
				allowed = append(allowed, fp)
			} else {
				blocked = append(blocked, fp)
			}
		}

		out[field] = FPolicies{
			Allowed: allowed,
			Blocked: blocked,
		}
	}

	return out
}