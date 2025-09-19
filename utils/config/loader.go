package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	middlewareTypes "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares/types"
	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
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
	BLOCKED_ENDPOINTS 	[]string 								`koanf:"blockedendpoints"`
	ALLOWED_ENDPOINTS 	[]string 								`koanf:"allowedendpoints"`
	VARIABLES 			map[string]any 							`koanf:"variables"`
	DATA_ALIASES 	map[string][]middlewareTypes.DataAlias 		`koanf:"dataaliases"`
	MESSAGE_TEMPLATE	string									`koanf:"messagetemplate"`
}

var ENV *ENV_ = &ENV_{
	CONFIG_PATH: os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR: os.Getenv("TOKENS_DIR"),
	API_TOKENS: []string{},
	SETTINGS: map[string]*SETTING_{

	},
	INSECURE: false,
}

func Load() {
	LoadDefaults()

	LoadConfig()

	LoadTokens()

	LoadEnv(userLayer)

	config = mergeLayers()

	normalizeKeys(config)
	templateConfig(config)

	InitTokens()

	InitEnv()

	log.Info("Finished Loading Configuration")

	log.Dev("Loaded Config:\n" + jsonutils.ToJson(config.All()))
	log.Dev("Loaded Token Configs:\n" + jsonutils.ToJson(tokensLayer.All()))
}

func InitEnv() {
	ENV.PORT = strconv.Itoa(config.Int("server.port"))

	ENV.LOG_LEVEL = strings.ToLower(config.String("loglevel"))
	
	ENV.API_URL = config.String("api.url")

	var settings SETTING_

	transformChildren(config, "settings.variables", transformVariables)

	config.Unmarshal("settings", &settings)

	ENV.SETTINGS["*"] = &settings
}

func LoadDefaults() {
	_, defErr := LoadFile(ENV.DEFAULTS_PATH, defaultsLayer, yaml.Parser())

	if defErr != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	_, conErr := LoadFile(ENV.CONFIG_PATH, userLayer, yaml.Parser())

	if conErr != nil {
		_, err := os.Stat(ENV.CONFIG_PATH)

		if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", conErr.Error())
		}
	}
}

func transformVariables(key string, value any) (string, any) {
	return strings.ToUpper(key), value
}