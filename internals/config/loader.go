package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/configutils"
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

var defaultsConf *configutils.Config
var userConf *configutils.Config
var tokenConf *configutils.Config

var mainConf *configutils.Config

func Load() {
	Clear()

	InitReload()

	LoadDefaults()

	LoadConfig()

	LoadTokens()

	userConf.LoadEnv()

	NormalizeConfig(defaultsConf)
	NormalizeConfig(userConf)
	
	NormalizeTokens()

	mainConf.MergeLayers(defaultsConf.Layer, userConf.Layer)

	mainConf.TemplateConfig()

	InitTokens()

	InitEnv()

	log.Info("Finished Loading Configuration")
}

func Log() {
	log.Dev("Loaded Config:", mainConf.Layer.All())
	log.Dev("Loaded Token Configs:", tokenConf.Layer.All())
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
	config.Load(data, "")
}

func NormalizeConfig(config *configutils.Config) {
	Normalize(config, "settings", &structure.SETTINGS{})
}

func Normalize(config *configutils.Config, path string, structure any) {
	data := config.Layer.Get(path)
	old, ok := data.(map[string]any)

	if !ok {
		log.Warn("Could not load `"+path+"`")
		return
	}

	// Create temporary config
	tmpConf := configutils.New()
	tmpConf.Load(old, "")
	
	// Apply transforms to the new config
	tmpConf.ApplyTransformFuncs(structure, "", transformFuncs)

	// Lowercase actual config
	LowercaseKeys(config)

	// Load temporary config back into paths
	config.Layer.Delete(path)
	
	config.Load(tmpConf.Layer.Get("").(map[string]any), path)
}

func InitReload() {
	reload := func() {
		Load()
		Log()
	}
	
	defaultsConf.OnLoad(reload)
	userConf.OnLoad(reload)
	tokenConf.OnLoad(reload)
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
