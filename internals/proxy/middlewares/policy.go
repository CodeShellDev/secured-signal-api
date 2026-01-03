package middlewares

import (
	"errors"
	"net/http"
	"reflect"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/requestkeys"
)

var Policy Middleware = Middleware{
	Name: "Policy",
	Use: policyHandler,
}

func policyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		policies := conf.SETTINGS.ACCESS.FIELD_POLICIES

		if policies == nil {
			policies = getConfig("").SETTINGS.ACCESS.FIELD_POLICIES
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
		}

		if body.Empty {
			body.Data = map[string]any{}
		}

		headerData := request.GetReqHeaders(req)

		shouldBlock, field := doBlock(body.Data, headerData, policies)

		if shouldBlock {
			logger.Warn("Client tried to use blocked field: ", field)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getPolicies(policies map[string]structure.FieldPolicy) (map[string]structure.FieldPolicy, map[string]structure.FieldPolicy) {
	blockedFields := map[string]structure.FieldPolicy{}
	allowedFields := map[string]structure.FieldPolicy{}

	for field, policy := range policies {
		switch policy.Action {
		case "block":
			blockedFields[field] = policy
		case "allow":
			allowedFields[field] = policy
		}
	}

	return allowedFields, blockedFields
}

func getField(key string, body map[string]any, headers map[string][]string) (any, error) {
	field := requestkeys.Parse(key)

	value := requestkeys.GetFromBodyAndHeaders(field, body, headers)

	if value != nil {
		return value, nil
	}

	return value, errors.New("field not found")
}

func doPoliciesApply(body map[string]any, headers map[string][]string, policies map[string]structure.FieldPolicy) (bool, string) {
	for key, policy := range policies {
		value, err := getField(key, body, headers)

		if err != nil {
			continue
		}

		switch asserted := value.(type) {
		case string:
			policyValue, ok := policy.Value.(string)

			if ok && asserted == policyValue {
				return true, key
			}
		case int:
			policyValue, ok := policy.Value.(int);

			if ok && asserted == policyValue {
				return true, key
			}
		case bool:
			policyValue, ok := policy.Value.(bool)

			if ok && asserted == policyValue {
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

func doBlock(body map[string]any, headers map[string][]string, policies map[string]structure.FieldPolicy) (bool, string) {
	if len(policies) == 0 {
		// default: allow all
		return false, ""
	}

	allowed, blocked := getPolicies(policies)

	var cause string

	isExplicitlyAllowed, cause := doPoliciesApply(body, headers, allowed)
	isExplicitlyBlocked, cause := doPoliciesApply(body, headers, blocked)
	
	// explicit allow > block
	if isExplicitlyAllowed {
		return false, cause
	}
	
	if isExplicitlyBlocked {
		return true, cause
	}

	// only allow policies -> block anything not allowed
	if len(allowed) > 0 && len(blocked) == 0 {
		return true, cause
	}

	// only block polcicies -> allow anything not blocked
	if len(blocked) > 0 && len(allowed) == 0 {
		return false, cause
	}

	// no match -> default: block all
	return true, cause
}
