package config

import (
	"reflect"
	"strconv"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/knadh/koanf/parsers/yaml"
)

func LoadTokens() {
	logger.Debug("Loading Configs in ", ENV.TOKENS_DIR)

	err := tokenConf.LoadDir("tokenconfigs", ENV.TOKENS_DIR, ".yml", yaml.Parser(), setTokenConfigName)

	if err != nil {
		logger.Error("Could not Load Configs in ", ENV.TOKENS_DIR, ": ", err.Error())
	}

	tokenConf.TemplateConfig()
}

func NormalizeTokens() {
	data := []map[string]any{}

	for _, config := range tokenConf.Layer.Slices("tokenconfigs") {
		tmpConf := configutils.New()
		tmpConf.Load(config.Raw(), "")

		Normalize("token", tmpConf, "", &structure.CONFIG{})
		
		data = append(data, tmpConf.Layer.Raw())
	}

	// Merge token configs together into new temporary config
	tokenConf.Load(data, "tokenconfigs")
}

func InitTokens() {
	apiTokens := DEFAULT.API.TOKENS

	for _, token := range apiTokens {
		ENV.CONFIGS[token] = DEFAULT
	}

	var tokenConfigs []structure.CONFIG

	tokenConf.Layer.Unmarshal("tokenconfigs", &tokenConfigs)

	config := parseTokenConfigs(tokenConfigs)

	for token, config := range config {
		apiTokens = append(apiTokens, token)

		ENV.CONFIGS[token] = &config
	}

	if len(apiTokens) <= 0 {
		logger.Warn("No API Tokens provided this is NOT recommended")

		logger.Info("Disabling Security Features due to incomplete Congfiguration")

		ENV.INSECURE = true

		// Set Blocked Endpoints on Config to User Layer Value
		// => effectively ignoring Default Layer
		DEFAULT.SETTINGS.ACCESS.ENDPOINTS = userConf.Layer.Strings("settings.access.endpoints")
	}

	if len(apiTokens) > 0 {
		logger.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")
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

func getSchemeTagByPointer(config any, tag string, fieldPointer any) string {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldValue := reflect.ValueOf(fieldPointer).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Addr().Interface() == fieldValue.Addr().Interface() {
			field := v.Type().Field(i)

			return field.Tag.Get(tag)
		}
	}

	return ""
}

func setTokenConfigName(config *configutils.Config, path string) {
	schema := reflect.TypeOf(structure.CONFIG{})

	nameField := getSchemeTagByPointer(schema, "koanf", structure.CONFIG{}.NAME)

	config.Layer.Set(nameField, path)
}