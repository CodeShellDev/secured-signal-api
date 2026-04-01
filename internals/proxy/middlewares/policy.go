package middlewares

import (
	"errors"
	"net/http"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
	c "github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
	"github.com/codeshelldev/secured-signal-api/utils/requestkeys"
)

var Policy Middleware = Middleware{
	Name: "Policy",
	Use: policyHandler,
}

func policyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		policies := conf.SETTINGS.ACCESS.FIELD_POLICIES.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.FIELD_POLICIES)

		if policies.Value == nil {
			next.ServeHTTP(w, req)
			return
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		body.EnsureNotNil()

		headers := request.GetReqHeaders(req)

		shouldBlock, field, err := isBlockedByPolicy(body.Data, headers, policies.Value.Compile())

		if err != nil {
			logger.Error("Could not perform policy checks: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return 
		}

		if shouldBlock {
			logger.Warn("Client tried to use blocked field: ", field)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getField(key string, body map[string]any, headers map[string][]string) (any, error) {
	field := requestkeys.Parse(key)

	value := requestkeys.GetFromBodyAndHeaders(field, body, headers)

	if value != nil {
		return value, nil
	}

	return nil, errors.New("field not found")
}

func doPoliciesApply(key string, body map[string]any, headers map[string][]string, policies []custom.FPolicy) (bool, string, error) {
	value, err := getField(key, body, headers)

	if err != nil {
		return false, "", nil
	}

	for _, policy := range policies {
		ok, err := policy.MatchRule.Match(value)

		if ok {
			return true, key, err
		}
	}

	return false, "", nil
}

func isBlockedByPolicy(body map[string]any, headers map[string][]string, p c.FieldPolicies) (bool, string, error) {
	policies := map[string]custom.FPolicies(p)

	if len(policies) == 0 {
		// default: allow all
		return false, "", nil
	}

	for field, policy := range policies {
		if len(policy.Allowed) == 0 && len(policy.Blocked) == 0 {
			continue
		}

		value, _ := getField(field, body, headers)

		if value == nil {
			continue
		}

		isExplicitlyAllowed, cause, err := doPoliciesApply(field, body, headers, policy.Allowed)

		if err != nil {
			return true, "", err
		}

		isExplicitlyBlocked, cause, err := doPoliciesApply(field, body, headers, policy.Blocked)

		if err != nil {
			return true, "", err
		}

		return checkBlockLogic(isExplicitlyAllowed, isExplicitlyBlocked, policy.Allowed, policy.Blocked), cause, nil
	}

	// default: allow all
	return false, "", nil
}
