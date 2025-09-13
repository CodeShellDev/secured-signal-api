package config

import (
	"strconv"

	middlewareTypes "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares/types"
	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/knadh/koanf/parsers/yaml"
)

var tokens []map[string]any

func LoadTokens() {
	log.Debug("Loading Configs ", ENV.TOKENS_DIR)

	LoadDir(ENV.TOKENS_DIR, yaml.Parser())

	log.Dev(utils.ToJson(tokens))
}

func InitTokens() {
	apiTokens := config.Strings("api.tokens")

	log.Dev(utils.ToJson(tokens))

	overrides := ParseTokenConfigs(tokens)

	for token, override := range overrides {
		apiTokens = append(apiTokens, token)

		ENV.SETTINGS[token] = &override
	}

	if len(apiTokens) <= 0 {
		log.Warn("No API TOKEN provided this is NOT recommended")

		log.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.INSECURE = true

		// Set Blocked Endpoints on Config to User Layer Value
		// => effectively ignoring Default Layer
		config.Set("blockedendpoints", userLayer.Strings("blockeendpoints"))
	}

	if len(apiTokens) > 0 {
		log.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")	

		ENV.API_TOKENS = apiTokens
	}
}

func ParseTokenConfigs(configs []map[string]any) (map[string]SETTING_) {
	settings := map[string]SETTING_{}

	for _, config := range configs {
		for _, token := range config["tokens"].([]string) {
			settings[token] = SETTING_{
				BLOCKED_ENDPOINTS: config["override.blockedendpoints"].([]string),
				VARIABLES: config["overrides.variables"].(map[string]any),
				MESSAGE_ALIASES: config["overrides.messagealiases"].([]middlewareTypes.MessageAlias),
			}
		}
	}

	return settings
}