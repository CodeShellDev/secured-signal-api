package runtime

import (
	"errors"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
)

func Validate() {
	ValidateConfig()
}

func ValidateConfig() {
	for _, conf := range config.ENV.CONFIGS {
		err, obj := ValidateEndpointRules(*conf)

		if err != nil {
			logger.Fatal("Error in endpoint rules: ", err.Error(), obj)
		}

		err, obj = ValidateFieldPolicyRules(*conf)

		if err != nil {
			logger.Fatal("Error in field policy rules: ", err.Error(), obj)
		}
	}
}

func ValidateEndpointRules(conf structure.CONFIG) (error, any) {
	if !conf.SETTINGS.ACCESS.ENDPOINTS.Set {
		return nil, nil
	}

	endpoints := conf.SETTINGS.ACCESS.ENDPOINTS.Value

	if endpoints == nil {
		return nil, nil
	}

	if len(endpoints.Allowed) != 0 {
		err := structure.StringMatchList(endpoints.Allowed).TestRules()

		if err != nil {
			return err, endpoints.Allowed
		}
	}

	if len(endpoints.Blocked) != 0 {
		err := structure.StringMatchList(endpoints.Blocked).TestRules()

		if err != nil {
			return err, endpoints.Blocked
		}
	}

	return nil, nil
}

func ValidateFieldPolicyRules(conf structure.CONFIG) (error, any) {
	if !conf.SETTINGS.ACCESS.FIELD_POLICIES.Set {
		return nil, nil
	}

	p := *conf.SETTINGS.ACCESS.FIELD_POLICIES.Value
	
	policies := map[string]custom.FPolicies(p.Compile())

	for field, policy := range policies {
		for _, item := range policy.Allowed {
			err := item.MatchRule.Test()

			if err != nil {
				return errors.New(field + ": " + err.Error()), item
			}
		}

		for _, item := range policy.Blocked {
			err := item.MatchRule.Test()

			if err != nil {
				return errors.New(field + ": " + err.Error()), item
			}
		}
	}

	return nil, nil
}