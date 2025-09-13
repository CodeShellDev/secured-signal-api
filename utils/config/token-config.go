package config

import (
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/knadh/koanf/parsers/yaml"
)

type TOKEN_CONFIG_ struct {
	TOKENS		[]string 	`koanf:"tokens"`
	OVERRIDES 	SETTING_	`koanf:"overrides"`
}

func LoadTokens() {
	log.Debug("Loading Configs ", ENV.TOKENS_DIR)

	LoadDir("tokenConfigs", ENV.TOKENS_DIR, tokensLayer, yaml.Parser())

	/*
	normalizeKeys(tokensLayer)

	templateConfig(tokensLayer)
	*/
}

func InitTokens() {
	apiTokens := config.Strings("api.tokens")

	var tokenConfigs []TOKEN_CONFIG_

	transformChildrenUnderArray(tokensLayer, "tokenConfigs", "overrides.variables", func(key string, value any) (string, any) {
		return strings.ToUpper(key), value
	})

	log.Dev(utils.ToJson(tokensLayer.All()))

	tokensLayer.Unmarshal("tokenConfigs", &tokenConfigs)

	log.Dev(utils.ToJson(tokenConfigs))

	overrides := parseTokenConfigs(tokenConfigs)

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

func parseTokenConfigs(configs []TOKEN_CONFIG_) (map[string]SETTING_) {
	settings := map[string]SETTING_{}

	for _, config := range configs {
		for _, token := range config.TOKENS {
			settings[token] = config.OVERRIDES
		}
	}

	return settings
}