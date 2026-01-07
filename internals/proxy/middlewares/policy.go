package middlewares

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"

	"github.com/codeshelldev/gotl/pkg/logger"
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
			return
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

func getPolicies(policies []structure.FieldPolicy) ([]structure.FieldPolicy, []structure.FieldPolicy) {
	blocked := []structure.FieldPolicy{}
	allowed := []structure.FieldPolicy{}

	for _, policy := range policies {
		switch policy.Action {
		case "block":
			blocked = append(blocked, policy)
		case "allow":
			allowed = append(allowed, policy)
		}
	}

	return allowed, blocked
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
				return re.MatchString(asserted), key
			}

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

func doBlock(body map[string]any, headers map[string][]string, policies map[string][]structure.FieldPolicy) (bool, string) {
	if len(policies) == 0 || policies == nil {
		// default: allow all
		return false, ""
	}

	var cause string

	for field, policy := range policies {
		value, _ := getField(field, body, headers)

		if value == nil {
			continue
		}

		allowed, blocked := getPolicies(policy)

		logger.Dev(allowed, blocked)

		isExplicitlyAllowed, cause := doPoliciesApply(field, body, headers, allowed)
		isExplicitlyBlocked, cause := doPoliciesApply(field, body, headers, blocked)

		logger.Dev(field, isExplicitlyAllowed, isExplicitlyBlocked)

		logger.Dev(policy)

		// block if explicitly blocked and no explicit allow exists
		if isExplicitlyBlocked && !isExplicitlyAllowed {
			return true, cause
		}
	}

	// default: allow all
	return false, cause
}
