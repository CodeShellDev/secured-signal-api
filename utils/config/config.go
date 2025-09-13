package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/codeshelldev/secured-signal-api/utils/safestrings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ENV_ struct {
	CONFIG_PATH			string
	DEFAULTS_PATH		string
	TOKENS_DIR			string
	LOG_LEVEL			string
	PORT 				string
	API_URL 			string
	API_TOKENS 			[]string
	INSECURE			bool
	BLOCKED_ENDPOINTS 	[]string
	VARIABLES 			map[string]any
	MESSAGE_ALIASES 	[]middlewares.MessageAlias
}

var ENV ENV_ = ENV_{
	CONFIG_PATH: os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR: os.Getenv("TOKENS_DIR"),
	API_TOKENS: []string{},
	BLOCKED_ENDPOINTS: []string{},
	MESSAGE_ALIASES: []middlewares.MessageAlias{},
	VARIABLES: map[string]any{},
	INSECURE: false,
}

var defaultsLayer = koanf.New(".")
var userLayer = koanf.New(".")

var config *koanf.Koanf

func InitEnv() {
	ENV.PORT = strconv.Itoa(config.Int("server.port"))

	ENV.LOG_LEVEL = config.String("loglevel")
	
	ENV.API_URL = config.String("api.url")

	apiTokens := config.Strings("api.tokens")

	if len(apiTokens) <= 0 {
		apiTokens = config.Strings("api.token")

		if len(apiTokens) <= 0 {
			log.Warn("No API TOKEN provided this is NOT recommended")

			log.Info("Disabling Security Features due to incomplete Congfiguration")

			ENV.INSECURE = true

			// Set Blocked Endpoints on Config to User Layer Value
			// => effectively ignoring Default Layer
			config.Set("blockedendpoints", userLayer.Strings("blockeendpoints"))
		}
	}

	if len(apiTokens) > 0 {
		log.Debug("Registered " + strconv.Itoa(len(apiTokens)) + " Tokens")	

		ENV.API_TOKENS = apiTokens
	}

	config.Unmarshal("messagealiases", &ENV.MESSAGE_ALIASES)

	transformChildren(config, "variables", func(key string, value any) (string, any) {
		return strings.ToUpper(key), value
	})

	config.Unmarshal("variables", &ENV.VARIABLES)

	ENV.BLOCKED_ENDPOINTS = config.Strings("blockedendpoints")
}

func Load() {
	log.Debug("Loading Config ", ENV.DEFAULTS_PATH)

	_, defErr := LoadFile(ENV.DEFAULTS_PATH, defaultsLayer, yaml.Parser())

	if defErr != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}

	log.Debug("Loading Config ", ENV.CONFIG_PATH)

	_, conErr := LoadFile(ENV.CONFIG_PATH, userLayer, yaml.Parser())

	if conErr != nil {
		_, err := os.Stat(ENV.CONFIG_PATH)

		if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", conErr.Error())
		}
	}

	log.Debug("Loading DotEnv")

	LoadEnv(userLayer)

	config = mergeLayers()

	normalizeKeys(config)

	templateConfig(config)

	InitEnv()

	log.Info("Finished Loading Configuration")

	log.Debug(utils.ToJson(config.All()))
	log.Debug(utils.ToJson(ENV))
}

func LoadFile(path string, config *koanf.Koanf, parser koanf.Parser) (koanf.Provider, error) {
	f := file.Provider(path)

	err := config.Load(f, parser)

	if err != nil {
		return nil, err
	}

	f.Watch(func(event any, err error) {
		if err != nil {
			return
		}

		log.Info("Config changed, Reloading...")

		Load()
	})

	return f, err
}

func templateConfig(config *koanf.Koanf) {
	data := config.All()

	for key, value := range data {
		str, isStr := value.(string)

		if isStr {
			templated := os.ExpandEnv(str)

			if templated != "" {
				data[key] = templated
			}
		}
	}

    config.Load(confmap.Provider(data, "."), nil)
}

func LoadEnv(config *koanf.Koanf) (koanf.Provider, error) {
	e := env.Provider(".", env.Opt{
		TransformFunc: normalizeEnv,
	})

	err := config.Load(e, nil)

	if err != nil {
		log.Fatal("Error loading env: ", err.Error())
	}

	return e, err
}

func mergeLayers() *koanf.Koanf {
	final := koanf.New(".")

	final.Merge(defaultsLayer)
	final.Merge(userLayer)

	return final
}

func normalizeKeys(config *koanf.Koanf) {
    data := map[string]any{}

    for _, key := range config.Keys() {
        lower := strings.ToLower(key)

        data[lower] = config.Get(key)
    }
    config.Load(confmap.Provider(data, "."), nil)
}

func transformChildren(config *koanf.Koanf, prefix string, transform func(key string, value any) (string, any)) error {
	var sub map[string]any
	if err := config.Unmarshal(prefix, &sub); err != nil {
		return err
	}

	transformed := make(map[string]any)
	for key, val := range sub {
		newKey, newVal := transform(key, val)

		transformed[newKey] = newVal
	}

	// Remove the old subtree by overwriting with empty map
	config.Load(confmap.Provider(map[string]any{
		prefix: map[string]any{},
	}, "."), nil)

	// Load the normalized subtree back in
	config.Load(confmap.Provider(map[string]any{
		prefix: transformed,
	}, "."), nil)

	return nil
}

func normalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", "")

	return key, safestrings.ToType(value)
}