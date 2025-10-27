package middlewares

import (
	"net/http"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/config/structure"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

var Policy Middleware = Middleware{
	Name: "Policy",
	Use: policyHandler,
}

func policyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		policies := settings.ACCESS.FIELD_POLOCIES

		if policies == nil {
			policies = getSettings("*").ACCESS.FIELD_POLOCIES
		}

		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
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

func hasField(field string, body map[string]any, headers map[string]any) bool {
	isHeader := strings.HasPrefix(field, "#")
	isBody := strings.HasPrefix(field, "@")

	fieldWithoutPrefix := field[:1]

	return (body[fieldWithoutPrefix] != nil && isBody) || (headers[fieldWithoutPrefix] != nil && isHeader)
}

func doBlock(body map[string]any, headers map[string]any, policies map[string]structure.FieldPolicy) (bool, string) {
	if policies == nil {
		return false, ""
	} else if len(policies) <= 0 {
		return false, ""
	}

	allowed, blocked := getPolicies(policies)

	var cause string

	var isExplictlyAllowed, isExplicitlyBlocked bool

	for field := range allowed {
		if hasField(field, body, headers) {
			isExplictlyAllowed = true
			cause = field
			break
		}
	}

	for field := range blocked {
		if hasField(field, body, headers) {
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
