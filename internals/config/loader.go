package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/configutils"
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

var defaultsConf = configutils.New()
var userConf = configutils.New()
var tokenConf = configutils.New()

var config = configutils.New()

func Load() {
	InitReload()

	LoadDefaults()

	LoadConfig()

	LoadTokens()

	userConf.LoadEnv()

	config.MergeLayers(defaultsConf.Layer, userConf.Layer)

	config.NormalizeKeys()
	config.TemplateConfig()

	InitTokens()

	InitEnv()

	log.Info("Finished Loading Configuration")

	log.Dev("Loaded Config:\n" + jsonutils.ToJson(config.Layer.All()))
	log.Dev("Loaded Token Configs:\n" + jsonutils.ToJson(tokenConf.Layer.All()))
}

func InitReload() {
	defaultsConf.OnLoad(Load)
	userConf.OnLoad(Load)
	tokenConf.OnLoad(Load)
}

func InitEnv() {
	ENV.PORT = strconv.Itoa(config.Layer.Int("service.port"))

	ENV.LOG_LEVEL = strings.ToLower(config.Layer.String("loglevel"))

	ENV.API_URL = config.Layer.String("api.url")

	var settings structure.SETTINGS

	config.TransformChildren("settings.message.variables", transformVariables)

	config.Layer.Unmarshal("settings", &settings)

	ENV.SETTINGS["*"] = &settings
}

func LoadDefaults() {
	_, err := defaultsConf.LoadFile(ENV.DEFAULTS_PATH, yaml.Parser())

	if err != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	_, err := userConf.LoadFile(ENV.CONFIG_PATH, yaml.Parser())

	if err != nil {
		_, fsErr := os.Stat(ENV.CONFIG_PATH)

		// Config File doesn't exist
		// => User is using Environment
		if errors.Is(fsErr, fs.ErrNotExist) {
			return
		}

		log.Error("Could not Load Config ", ENV.CONFIG_PATH, ": ", err.Error())
	}
}

func transformVariables(key string, value any) (string, any) {
	return strings.ToUpper(key), value
}
