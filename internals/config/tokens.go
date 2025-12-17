package config

import (
	"strconv"

	"github.com/codeshelldev/gotl/pkg/configutils"
	log "github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/knadh/koanf/parsers/yaml"
)

func LoadTokens() {
	log.Debug("Loading Configs in ", ENV.TOKENS_DIR)

	err := tokenConf.LoadDir("tokenconfigs", ENV.TOKENS_DIR, ".yml", yaml.Parser())

	if err != nil {
		log.Error("Could not Load Configs in ", ENV.TOKENS_DIR, ": ", err.Error())
	}

	tokenConf.TemplateConfig()
}

func NormalizeTokens() {
	data := []map[string]any{}

	for _, config := range tokenConf.Layer.Slices("tokenconfigs") {
		tmpConf := configutils.New()
		tmpConf.Load(config.Raw(), "")

		Normalize("token", tmpConf, "", &structure.SETTINGS{})
		
		data = append(data, tmpConf.Layer.Raw())
	}

	// Merge token configs together into new temporary config
	tokenConf.Load(data, "tokenconfigs")
}

func InitTokens() {
	apiTokens := DEFAULT.API.TOKENS

	var tokenConfigs []structure.CONFIG

	tokenConf.Layer.Unmarshal("tokenconfigs", &tokenConfigs)

	log.Dev("TokenConfigs:", tokenConfigs)

	config := parseTokenConfigs(tokenConfigs)

	for token, config := range config {
		apiTokens = append(apiTokens, token)

		ENV.CONFIGS[token] = &config
	}

	if len(apiTokens) <= 0 {
		log.Warn("No API Tokens provided this is NOT recommended")

		log.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.INSECURE = true

		// Set Blocked Endpoints on Config to User Layer Value
		// => effectively ignoring Default Layer
		DEFAULT.SETTINGS.ACCESS.ENDPOINTS = userConf.Layer.Strings("settings.access.endpoints")
	}

	if len(apiTokens) > 0 {
		log.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")
	}
}

func parseTokenConfigs(configArray []structure.CONFIG) map[string]structure.CONFIG {
	configs := map[string]structure.CONFIG{}

	for _, config := range configArray {
		for _, token := range config.API.TOKENS {
			configs[token] = config
		}
	}

	return configs
}
