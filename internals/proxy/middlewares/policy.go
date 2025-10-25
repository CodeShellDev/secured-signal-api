package middlewares

import (
	"net/http"
	"slices"
	"strings"

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

func getPolicies(policies []string) ([]string, []string) {
	blockedFields := []string{}
	allowedFields := []string{}

	for _, field := range policies {
		field, block := strings.CutPrefix(field, "!")

		if block {
			blockedFields = append(blockedFields, field)
		} else {
			allowedFields = append(allowedFields, field)
		}
	}

	return allowedFields, blockedFields
}

func doBlock(body map[string]any, headers map[string]any, policies []string) (bool, string) {
	if policies == nil {
		return false, ""
	} else if len(policies) <= 0 {
		return false, ""
	}

	allowed, blocked := getPolicies(policies)

	var blockField string

	isExplicitlyBlocked := slices.ContainsFunc(blocked, func(try string) bool {
		isHeader := strings.HasPrefix(try, "#")
		isBody := strings.HasPrefix(try, "@")

		if body[try] != nil && isBody {
			blockField = try
			return true
		}

		if headers[try] != nil && isHeader {
			blockField = try
			return true
		}

		return false
	})

	isExplictlyAllowed := slices.ContainsFunc(allowed, func(try string) bool {
		isHeader := strings.HasPrefix(try, "#")
		isBody := strings.HasPrefix(try, "@")

		if body[try] != nil && isBody {
			return true
		}

		if headers[try] != nil && isHeader {
			return true
		}

		return false
	})

	// Block all except explicitly Allowed
	if len(blocked) == 0 && len(allowed) != 0 {
		return !isExplictlyAllowed, blockField
	}

	// Allow all except explicitly Blocked
	if len(allowed) == 0 && len(blocked) != 0 {
		return isExplicitlyBlocked, blockField
	}

	// Excplicitly Blocked except excplictly Allowed
	if len(blocked) != 0 && len(allowed) != 0 {
		return isExplicitlyBlocked && !isExplictlyAllowed, blockField
	}

	// Block all
	return true, ""
}
