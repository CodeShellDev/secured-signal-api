package config

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/knadh/koanf/parsers/yaml"
)

const tokenConfigsPath = "tokenconfigs"

func LoadTokens() {
	logger.Debug("Loading Configs in ", ENV.TOKENS_DIR)

	err := tokenConf.LoadDir(tokenConfigsPath, ENV.TOKENS_DIR, ".yml", yaml.Parser(), setTokenConfigName)

	if err != nil {
		logger.Error("Could not Load Configs in ", ENV.TOKENS_DIR, ": ", err.Error())
	}
}

func NormalizeTokens() {
	data := []any{}

	for _, config := range tokenConf.Layer.Slices(tokenConfigsPath) {
		tmpConf := configutils.New()
		tmpConf.Load(config.Raw(), "")

		Normalize("token", tmpConf, "", &structure.CONFIG{})
		
		data = append(data, tmpConf.Layer.Raw())
	}

	// Merge token configs together into new temporary config
	tokenConf.Layer.Set(tokenConfigsPath, data)
}

func InitTokens() {
	apiTokens := parseAuthTokens(*DEFAULT)

	for _, token := range apiTokens {
		ENV.CONFIGS[token] = DEFAULT
	}

	configs := parseTokenConfigs(tokenConf)

	for token, config := range configs {
		apiTokens = append(apiTokens, token)

		config.TYPE = structure.TOKEN

		ENV.CONFIGS[token] = &config
	}

	if len(apiTokens) <= 0 {
		logger.Warn("No API Tokens provided this is NOT recommended")

		logger.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.INSECURE = true

		// Set Blocked Endpoints on Config to User Layer Value
		// => effectively ignoring Default Layer
		userConf.Unmarshal("settings.access.endpoints", &DEFAULT.SETTINGS.ACCESS.ENDPOINTS.Value)
	}

	if len(apiTokens) > 0 {
		logger.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")
	}

	ENV.TOKENS = apiTokens
}

func parseTokenConfigs(config *configutils.Config) map[string]structure.CONFIG {
	configs := map[string]structure.CONFIG{}

	for _, c := range config.Layer.Slices(tokenConfigsPath) {
		tmpConf := configutils.New()
		tmpConf.Load(c.Raw(), "")

		templateConfigWithVariables(tmpConf)
		
		var configData structure.CONFIG

		tmpConf.Unmarshal("", &configData)

		tokens := parseAuthTokens(configData)

		for _, token := range tokens {
			configs[token] = configData
		}
	}

	return configs
}

func parseAuthTokens(config structure.CONFIG) []string {
	tokens := config.API.TOKENS

	for _, token := range config.API.AUTH.TOKENS {
		tokens = append(tokens, token.Set...)
	}

	return tokens
}

func setTokenConfigName(config *configutils.Config, p string) {
	schema := structure.CONFIG{
		NAME: "",
	}

	nameField := configutils.GetSchemeTagByFieldPointer(&schema, "koanf", &schema.NAME)

	filename := filepath.Base(p)
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	name := config.Layer.String(nameField)

	if strings.TrimSpace(name) == "" {
		config.Layer.Set(nameField, filenameWithoutExt)
	}
}