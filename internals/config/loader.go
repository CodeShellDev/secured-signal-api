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
	"github.com/knadh/koanf/providers/confmap"
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

var mainConf = configutils.New()

func Load() {
	Clear()

	InitReload()

	LoadDefaults()

	LoadConfig()

	LoadTokens()

	userConf.LoadEnv()

	mainConf.MergeLayers(defaultsConf.Layer, userConf.Layer)

	NormalizeConfig()
	NormalizeTokens()

	mainConf.TemplateConfig()

	InitTokens()

	InitEnv()

	log.Info("Finished Loading Configuration")

	log.Dev("Loaded Config:\n" + jsonutils.ToJson(mainConf.Layer.All()))
	log.Dev("Loaded Token Configs:\n" + jsonutils.ToJson(tokenConf.Layer.All()))
}

func Clear() {
	defaultsConf = configutils.New()
	userConf = configutils.New()
	tokenConf = configutils.New()
	mainConf = configutils.New()
}

func LowercaseKeys(config *configutils.Config) {
	data := map[string]any{}

	for _, key := range config.Layer.Keys() {
		lower := strings.ToLower(key)

		data[lower] = config.Layer.Get(key)
	}

	config.Layer.Delete("")
	config.Layer.Load(confmap.Provider(data, "."), nil)
}

func NormalizeConfig() {
	settings := mainConf.Layer.Get("settings")
	old, ok := settings.(map[string]any)

	if !ok {
		log.Warn("Could not load `settings`")
		return
	}

	// Create temporary configs
	tmpConf := configutils.New()
	tmpConf.Layer.Load(confmap.Provider(old, "."), nil)
	
	// Apply transforms to the new configs
	tmpConf.ApplyTransformFuncs(&structure.SETTINGS{}, "", transformFuncs)

	// Lowercase actual configs
	LowercaseKeys(mainConf)

	// Load temporary configs back into paths
	mainConf.Layer.Delete("settings")

	log.Dev("Loading:\n--------------------------------------\n", jsonutils.ToJson(mainConf.Layer.All()), "\n--------------------------------------")

	mainConf.Layer.Load(confmap.Provider(tmpConf.Layer.All(), "settings"), nil)
}

func InitReload() {
	defaultsConf.OnLoad(Load)
	userConf.OnLoad(Load)
	tokenConf.OnLoad(Load)
}

func InitEnv() {
	ENV.PORT = strconv.Itoa(mainConf.Layer.Int("service.port"))

	ENV.LOG_LEVEL = strings.ToLower(mainConf.Layer.String("loglevel"))

	ENV.API_URL = mainConf.Layer.String("api.url")

	var settings structure.SETTINGS

	mainConf.Layer.Unmarshal("settings", &settings)

	ENV.SETTINGS["*"] = &settings
}

func LoadDefaults() {
	_, err := defaultsConf.LoadFile(ENV.DEFAULTS_PATH, yaml.Parser())

	log.Dev("Defaults:\n--------------------------------------\n", jsonutils.ToJson(defaultsConf.Layer.All()), "\n--------------------------------------")

	if err != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	_, err := userConf.LoadFile(ENV.CONFIG_PATH, yaml.Parser())

	log.Dev("User:\n--------------------------------------\n", jsonutils.ToJson(userConf.Layer.All()), "\n--------------------------------------")

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