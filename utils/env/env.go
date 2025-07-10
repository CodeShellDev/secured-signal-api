package env

import (
	"encoding/json"
	"os"

	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

type ENV_ struct {
	PORT 				string
	API_URL 			string
	API_TOKEN 			string
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
		"NUMBER": os.Getenv("SENDER"),
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

	ENV.API_TOKEN = os.Getenv("API_TOKEN")

	blockedEndpointJSON := os.Getenv("BLOCKED_ENDPOINTS")
	recipientsJSON := os.Getenv("DEFAULT_RECIPIENTS")
	messageAliasesJSON := os.Getenv("MESSAGE_ALIASES")
	variablesJSON := os.Getenv("VARIABLES")

	log.Info("Loaded Environment Variables")

	if ENV.API_TOKEN == "" {
		log.Warn("No API TOKEN provided this is NOT recommended")

		log.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.BLOCKED_ENDPOINTS = []string{}
	}

	if blockedEndpointJSON != "" {
		var blockedEndpoints []string

		err := json.Unmarshal([]byte(blockedEndpointJSON), &blockedEndpoints)

		if err != nil {
			log.Error("Could not decode Blocked Endpoints: ", blockedEndpointJSON)
		}

		ENV.BLOCKED_ENDPOINTS = blockedEndpoints
	}

	if messageAliasesJSON != "" {
		var msgAliases []middlewares.MessageAlias

		err := json.Unmarshal([]byte(messageAliasesJSON), &msgAliases)

		if err != nil {
			log.Error("Could not decode Message Aliases ", variablesJSON)
		}

		ENV.MESSAGE_ALIASES = msgAliases
	}

	if variablesJSON != "" {
		var variables map[string]interface{}

		err := json.Unmarshal([]byte(variablesJSON), &variables)

		if err != nil {
			log.Error("Could not decode Variables ", variablesJSON)
		}

		ENV.VARIABLES = variables
	}

	if recipientsJSON != "" {
		var recipients []string

		err := json.Unmarshal([]byte(recipientsJSON), &recipients)

		if err != nil {
			log.Error("Could not decode Variables ", variablesJSON)
		}

		ENV.VARIABLES["RECIPIENTS"] = recipients
	}
}