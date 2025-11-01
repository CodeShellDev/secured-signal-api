package middlewares

import (
	"errors"
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
	"github.com/codeshelldev/secured-signal-api/utils/request/requestkeys"
)

var Policy Middleware = Middleware{
	Name: "Policy",
	Use: policyHandler,
}

func policyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		policies := settings.ACCESS.FIELD_POLICIES

		if policies == nil {
			policies = getSettings("*").ACCESS.FIELD_POLICIES
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
		}

		if body.Empty {
			body.Data = map[string]any{}
		}

		headerData := request.GetReqHeaders(req)

		shouldBlock, field := doBlock(body.Data, headerData, policies)

		if shouldBlock {
			log.Warn("User tried to use blocked field: ", field)
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

func doBlock(body map[string]any, headers map[string][]string, policies map[string]structure.FieldPolicy) (bool, string) {
	if policies == nil {
		return false, ""
	} else if len(policies) <= 0 {
		return false, ""
	}

	allowed, blocked := getPolicies(policies)

	var cause string

	var isExplictlyAllowed, isExplicitlyBlocked bool

	for field, policy := range allowed {
		value, err := getField(field, body, headers)

		if value == policy.Value && err == nil {
			isExplictlyAllowed = true
			cause = field
			break
		}
	}

	for field, policy := range blocked {
		value, err := getField(field, body, headers)

		if value == policy.Value && err == nil {
			isExplicitlyBlocked = true
			cause = field
			break
		}
	}

	// Block all except explicitly Allowed
	if len(blocked) == 0 && len(allowed) != 0 {
		return !isExplictlyAllowed, cause
	}

	// Allow all except explicitly Blocked
	if len(allowed) == 0 && len(blocked) != 0 {
		return isExplicitlyBlocked, cause
	}

	// Excplicitly Blocked except excplictly Allowed
	if len(blocked) != 0 && len(allowed) != 0 {
		return isExplicitlyBlocked && !isExplictlyAllowed, cause
	}

	// Block all
	return true, ""
}
