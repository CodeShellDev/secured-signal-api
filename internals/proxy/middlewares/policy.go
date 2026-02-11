package middlewares

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
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

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		if body.Empty {
			body.Data = map[string]any{}
		}

		headerData := request.GetReqHeaders(req)

		shouldBlock, field := isBlockedByPolicy(body.Data, headerData, policies)

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

func doPoliciesApply(key string, body map[string]any, headers map[string][]string, policies []structure.FieldPolicy) (bool, string) {
	value, err := getField(key, body, headers)

	if err != nil {
		return false, ""
	}

	for _, policy := range policies {
		switch asserted := value.(type) {
		case string:
			policyValue, ok := policy.Value.(string)

			re, err := regexp.Compile(policyValue)

			if err == nil {
				if re.MatchString(asserted) {
					return true, key
				}
				continue
			}

			if ok && asserted == policyValue {
				return true, key
			}
		case int:
			policyValue, ok := policy.Value.(int)

			if ok && asserted == policyValue {
				return true, key
			}
		case float64:
			var policyValue float64

			// needed for json
			switch assertedValue := policy.Value.(type) {
			case int:
				policyValue = float64(assertedValue)
			case float64:
				policyValue = assertedValue
			default:
				continue
			}

			if asserted == policyValue {
				return true, key
			}
		default:
			if reflect.DeepEqual(value, policy.Value) {
				return true, key
			}
		}
	}

	return false, ""
}

func isBlockedByPolicy(body map[string]any, headers map[string][]string, policies map[string]structure.FieldPolicies) (bool, string) {
	if len(policies) == 0 || policies == nil {
		// default: allow all
		return false, ""
	}

	for field, policy := range policies {
		if len(policy.Allow) == 0 || len(policy.Block) == 0 {
			continue
		}

		value, _ := getField(field, body, headers)

		if value == nil {
			continue
		}

		isExplicitlyAllowed, cause := doPoliciesApply(field, body, headers, policy.Allow)
		isExplicitlyBlocked, cause := doPoliciesApply(field, body, headers, policy.Block)

		// explicit allow > block
		if isExplicitlyAllowed {
			return false, cause
		}
		
		if isExplicitlyBlocked {
			return true, cause
		}

		// allow rules -> default deny
		if len(policy.Allow) > 0 {
			return true, cause
		}
		
		// only block rules -> default allow
		if len(policy.Block) > 0 {
			return false, cause
		}

		// safety net -> block
		return true, "safety net"
	}

	// default: allow all
	return false, ""
}
