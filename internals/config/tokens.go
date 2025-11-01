package config

import (
	"strconv"

	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/configutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/knadh/koanf/parsers/yaml"
)

type TOKEN_CONFIG_ struct {
	TOKENS    []string 				`koanf:"tokens"`
	OVERRIDES structure.SETTINGS 	`koanf:"overrides"`
}

func LoadTokens() {
	log.Debug("Loading Configs in ", ENV.TOKENS_DIR)

	err := tokenConf.LoadDir("tokenconfigs", ENV.TOKENS_DIR, ".yml", yaml.Parser())

	if err != nil {
		log.Error("Could not Load Configs in ", ENV.TOKENS_DIR, ": ", err.Error())
	}

	tokenConf.TemplateConfig()
}

func NormalizeTokens() {
	configArray := []map[string]any{}

	for _, config := range tokenConf.Layer.Slices("tokenconfigs") {
		tmpConf := configutils.New()
		tmpConf.Load(config.All(), "")

		Normalize(tmpConf, "overrides", &structure.SETTINGS{})

		tmpConf.ApplyTransformFuncs(&structure.SETTINGS{}, "overrides", transformFuncs)

		configArray = append(configArray, tmpConf.Layer.All())
	}

	// Merge token configs together into new temporary config
	tokenConf.Layer.Set("tokenconfigs", configArray)
}

func InitTokens() {
	apiTokens := mainConf.Layer.Strings("api.tokens")

	var tokenConfigs []TOKEN_CONFIG_

	tokenConf.Layer.Unmarshal("tokenconfigs", &tokenConfigs)

	overrides := parseTokenConfigs(tokenConfigs)

	for token, override := range overrides {
		apiTokens = append(apiTokens, token)

		ENV.SETTINGS[token] = &override
	}

	if len(apiTokens) <= 0 {
		log.Warn("No API Tokens provided this is NOT recommended")

		log.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.INSECURE = true

		// Set Blocked Endpoints on Config to User Layer Value
		// => effectively ignoring Default Layer
		mainConf.Layer.Set("settings.access.endpoints", userConf.Layer.Strings("settings.access.endpoints"))
	}

	if len(apiTokens) > 0 {
		log.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")

		ENV.API_TOKENS = apiTokens
	}
}

func parseTokenConfigs(configs []TOKEN_CONFIG_) map[string]structure.SETTINGS {
	settings := map[string]structure.SETTINGS{}

	for _, config := range configs {
		for _, token := range config.TOKENS {
			settings[token] = config.OVERRIDES
		}
	}

	return settings
}
