package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	middlewareTypes "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares/types"
	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/knadh/koanf/parsers/yaml"
)

type ENV_ struct {
	CONFIG_PATH			string
	DEFAULTS_PATH		string
	TOKENS_DIR			string
	LOG_LEVEL			string
	PORT 				string
	API_URL 			string
	API_TOKENS 			[]string
	SETTINGS			map[string]*SETTING_
	INSECURE			bool
}

type SETTING_ struct {
	BLOCKED_ENDPOINTS 	[]string 						`koanf:"blockedendpoints"`
	ALLOWED_ENDPOINTS 	[]string 						`koanf:"allowedendpoints"`
	VARIABLES 			map[string]any 					`koanf:"variables"`
	MESSAGE_ALIASES 	[]middlewareTypes.MessageAlias 	`koanf:"messagealiases"`
}

var ENV *ENV_ = &ENV_{
	CONFIG_PATH: os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR: os.Getenv("TOKENS_DIR"),
	API_TOKENS: []string{},
	SETTINGS: map[string]*SETTING_{
		"*": {
			BLOCKED_ENDPOINTS: []string{},
			ALLOWED_ENDPOINTS: []string{},
			MESSAGE_ALIASES: []middlewareTypes.MessageAlias{},
			VARIABLES: map[string]any{},
		},
	},
	INSECURE: false,
}

func InitEnv() {
	ENV.PORT = strconv.Itoa(config.Int("server.port"))

	ENV.LOG_LEVEL = config.String("loglevel")
	
	ENV.API_URL = config.String("api.url")

	transformChildren(config, "variables", func(key string, value any) (string, any) {
		return strings.ToUpper(key), value
	})

	config.Unmarshal("settings", ENV.SETTINGS["*"])
}

func Load() {
	LoadDefaults()

	LoadConfig()

	LoadTokens()

	log.Debug("Loading DotEnv")

	LoadEnv(userLayer)

	config = mergeLayers()

	normalizeKeys(config)
	templateConfig(config)

	InitTokens()

	InitEnv()

	log.Info("Finished Loading Configuration")

	log.Dev("Loaded Config:\n" + utils.ToJson(config.All()))
	log.Dev("Loaded Token Configs:\n" + utils.ToJson(tokensLayer.All()))
}

func LoadDefaults() {
	log.Debug("Loading Config ", ENV.DEFAULTS_PATH)

	_, defErr := LoadFile(ENV.DEFAULTS_PATH, defaultsLayer, yaml.Parser())

	if defErr != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	log.Debug("Loading Config ", ENV.CONFIG_PATH)

	_, conErr := LoadFile(ENV.CONFIG_PATH, userLayer, yaml.Parser())

	if conErr != nil {
		_, err := os.Stat(ENV.CONFIG_PATH)

		if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", conErr.Error())
		}
	}
}