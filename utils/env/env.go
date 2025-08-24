package env

import (
	"os"
	"strconv"

	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type ENV_ struct {
	PORT 				string
	API_URL 			string
	API_TOKENS 			[]string
	BLOCKED_ENDPOINTS 	[]string
	VARIABLES 			map[string]any
	MESSAGE_ALIASES 	[]middlewares.MessageAlias
}

var ENV ENV_ = ENV_{
	BLOCKED_ENDPOINTS: []string {
		"/v1/about",
		"/v1/configuration",
		"/v1/devices",
		"/v1/register",
		"/v1/unregister",
		"/v1/qrcodelink",
		"/v1/accounts",
		"/v1/contacts",
	},
	VARIABLES: map[string]any {
		"RECIPIENTS": []string{},
		"NUMBER": os.Getenv("NUMBER"),
	},
	MESSAGE_ALIASES: []middlewares.MessageAlias{
		{
			Alias:    "msg",
			Priority: 100,
		},
		{
			Alias:    "content",
			Priority: 99,
		},
		{
			Alias:    "description",
			Priority: 98,
		},
		{
			Alias:    "text",
			Priority: 20,
		},
		{
			Alias:    "body",
			Priority: 15,
		},
		{
			Alias:    "summary",
			Priority: 10,
		},
		{
			Alias:    "details",
			Priority: 9,
		},
		{
			Alias:    "payload",
			Priority: 2,
		},
		{
			Alias:    "data",
			Priority: 1,
		},
	},
}

func Load() {
	ENV.PORT = os.Getenv("PORT")
	ENV.API_URL = os.Getenv("SIGNAL_API_URL")

	apiToken := os.Getenv("API_TOKENS")

	if apiToken == "" {
		apiToken = os.Getenv("API_TOKEN")
	}

	blockedEndpointJSON := os.Getenv("BLOCKED_ENDPOINTS")
	recipientsJSON := os.Getenv("RECIPIENTS")
	messageAliasesJSON := os.Getenv("MESSAGE_ALIASES")
	variablesJSON := os.Getenv("VARIABLES")

	log.Info("Loaded Environment Variables")

	apiTokens, err := utils.StringToArray(apiToken)

	if err != nil {
		log.Warn("No API TOKEN provided this is NOT recommended")

		log.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.BLOCKED_ENDPOINTS = []string{}
	} else {
		log.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")

		ENV.API_TOKENS = apiTokens
	}

	if blockedEndpointJSON != "" {
		ENV.BLOCKED_ENDPOINTS = utils.GetJson[[]string](blockedEndpointJSON)
	}

	if messageAliasesJSON != "" {
		ENV.MESSAGE_ALIASES = utils.GetJson[[]middlewares.MessageAlias](messageAliasesJSON)
	}

	if variablesJSON != "" {
		ENV.VARIABLES = utils.GetJson[map[string]any](variablesJSON)
	}

	if recipientsJSON != "" {
		ENV.VARIABLES["RECIPIENTS"] = utils.GetJson[[]string](recipientsJSON)
	}
}