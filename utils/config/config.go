package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	utils "github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ENV_ struct {
	CONFIG_PATH			string
	DEFAULTS_PATH		string
	TOKENS_DIR			string
	PORT 				string
	API_URL 			string
	API_TOKENS 			[]string
	BLOCKED_ENDPOINTS 	[]string
	VARIABLES 			map[string]any
	MESSAGE_ALIASES 	[]middlewares.MessageAlias
}

var ENV ENV_ = ENV_{
	CONFIG_PATH: os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR: os.Getenv("TOKENS_DIR"),
}

var config = koanf.New(".")

func LoadIntoENV() {
	ENV.PORT = strconv.Itoa(config.Int("server.port"))
	
	ENV.API_URL = config.String("api.url")

	apiTokens := config.Strings("api.tokens")

	if len(apiTokens) <= 0 {
		apiTokens = config.Strings("api.token")
	}

	ENV.API_TOKENS = apiTokens

	ENV.BLOCKED_ENDPOINTS = config.Strings("blockedendpoints")

	ENV.VARIABLES = config.Get("variables").(map[string]any)
	ENV.MESSAGE_ALIASES = config.Get("messagealiases").([]middlewares.MessageAlias)

	ENV.VARIABLES["NUMBER"] = config.String("number")
	ENV.VARIABLES["RECIPIENTS"] = config.Strings("recipients")
}

func Load() {
	log.Debug("Loading Default Config ", ENV.DEFAULTS_PATH)

	defErr := LoadFile(ENV.DEFAULTS_PATH, yaml.Parser())

	if defErr != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}

	log.Debug("Loading Config ", ENV.CONFIG_PATH)

	conErr := LoadFile(ENV.CONFIG_PATH, yaml.Parser())

	if conErr != nil {
		_, err := os.Stat(ENV.CONFIG_PATH)

		if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", conErr.Error())
		}
	}

	log.Debug("Loading DotEnv")
	LoadDotEnv()

	normalizeKeys()

	LoadIntoENV()

	log.Info("Finished Loading Configuration")
}

func LoadFile(path string, parser koanf.Parser) error {
	f := file.Provider(path)

	err := config.Load(f, parser)

	if err != nil {
		return err
	}

	f.Watch(func(event any, err error) {
		if err != nil {
			return
		}

		log.Info("Config changed, Reloading...")

		Load()
	})

	return err
}

func LoadDotEnv() error {
	e := env.ProviderWithValue("", ".", normalizeEnv)

	err := config.Load(e, dotenv.Parser())

	if err != nil {
		log.Fatal("Error loading env: ", err.Error())
	}

	return err
}

func normalizeKeys() {
    data := map[string]any{}

    for _, key := range config.Keys() {
        lower := strings.ToLower(key)

        data[lower] = config.Get(key)
    }
    config.Load(confmap.Provider(data, "."), nil)
}

func normalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(strings.ReplaceAll(key, "__", "."))

	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		data, err := utils.GetJsonSafe[any](value)

		if data != nil && err == nil {
			return key, data
		}
	}

	if strings.Contains(value, ",") {
		items := utils.StringToArray(value)
		
		return key, items
	}

	intValue, intErr := strconv.Atoi(value)

	if intErr == nil {
		return key, intValue
	}

	return key, value
}