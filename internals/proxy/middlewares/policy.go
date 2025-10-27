package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/jsonutils"
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

func getField(field string, body map[string]any, headers map[string]any) (any, error) {
	isHeader := strings.HasPrefix(field, "#")
	isBody := strings.HasPrefix(field, "@")

	fieldWithoutPrefix := field[:1]

	var value any

	if body[fieldWithoutPrefix] != nil && isBody {
		value = body[fieldWithoutPrefix]
	} else if headers[fieldWithoutPrefix] != nil && isHeader {
		value = headers[fieldWithoutPrefix]
	}

	if value != nil {
		return value, nil
	}

	return value, errors.New("field not found")
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

	for field, policy := range allowed {
		value, err := getField(field, body, headers)

		log.Dev("Checking ", field, "...")
		log.Dev("Got Value of ", jsonutils.ToJson(value))

		if value == policy.Value && err == nil {
			isExplictlyAllowed = true
			cause = field
			break
		}
	}

	for field, policy := range blocked {
		value, err := getField(field, body, headers)

		log.Dev("Checking ", field, "...")
		log.Dev("Got Value of ", jsonutils.ToJson(value))

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
