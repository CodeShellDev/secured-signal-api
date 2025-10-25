package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/config/structure"
	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"

	"github.com/knadh/koanf/parsers/yaml"
)

var ENV *structure.ENV = &structure.ENV{
	CONFIG_PATH:   os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR:    os.Getenv("TOKENS_DIR"),
	FAVICON_PATH:  os.Getenv("FAVICON_PATH"),
	API_TOKENS:    []string{},
	SETTINGS:      map[string]*structure.SETTINGS{},
	INSECURE:      false,
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
	ENV.PORT = strconv.Itoa(config.Int("service.port"))

	ENV.LOG_LEVEL = strings.ToLower(config.String("loglevel"))

	ENV.API_URL = config.String("api.url")

	var settings structure.SETTINGS

	transformChildren(config, "settings.message.variables", transformVariables)

	config.Unmarshal("settings", &settings)

	ENV.SETTINGS["*"] = &settings
}

func LoadDefaults() {
	_, err := LoadFile(ENV.DEFAULTS_PATH, defaultsLayer, yaml.Parser())

	if err != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	_, err := LoadFile(ENV.CONFIG_PATH, userLayer, yaml.Parser())

	if err != nil {
		_, fsErr := os.Stat(ENV.CONFIG_PATH)

		// Config File doesn't exist
		// => User is using Environment
		if errors.Is(fsErr, fs.ErrNotExist) {
			return
		}
	}

	log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", err.Error())
}

func transformVariables(key string, value any) (string, any) {
	return strings.ToUpper(key), value
}
