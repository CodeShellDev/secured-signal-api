package config

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/codeshelldev/gotl/pkg/configutils"
	log "github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/stringutils"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"

	"github.com/knadh/koanf/parsers/yaml"
)

var ENV *structure.ENV = &structure.ENV{
	CONFIG_PATH:   	os.Getenv("CONFIG_PATH"),
	DEFAULTS_PATH: 	os.Getenv("DEFAULTS_PATH"),
	TOKENS_DIR:    	os.Getenv("TOKENS_DIR"),
	FAVICON_PATH:  	os.Getenv("FAVICON_PATH"),
	INSECURE:      	false,

	CONFIGS:       	map[string]*structure.CONFIG{},
}

var DEFAULT	*structure.CONFIG

var defaultsConf *configutils.Config
var userConf *configutils.Config
var envConf *configutils.Config
var tokenConf *configutils.Config

var mainConf *configutils.Config

func Load() {
	Clear()

	InitReload()

	LoadDefaults()

	LoadConfig()

	LoadTokens()

	//NormalizeConfig("", defaultsConf)
	//NormalizeConfig("config", userConf)

	envConf.LoadEnv(normalizeEnv)

	NormalizeConfig("env", envConf)

	userConf.MergeLayers(envConf.Layer)
	
	mainConf.MergeLayers(defaultsConf.Layer, userConf.Layer)

	mainConf.TemplateConfig()

	//NormalizeTokens()

	InitConfig()

	InitTokens()

	log.Info("Finished Loading Configuration")
}

func Log() {
	// TODO: Change back to `log.Dev()` as soon as config parsing is working again.
	log.Info("Loaded Config:", mainConf.Layer.All())
	log.Info("Loaded Token Configs:", tokenConf.Layer.All())
	log.Info("Parsed Configs: ", ENV)
}

func Clear() {
	defaultsConf = configutils.New()
	userConf = configutils.New()
	envConf = configutils.New()
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

func NormalizeConfig(id string, config *configutils.Config) {
	Normalize(id, config, "", &structure.CONFIG{})
}

func Normalize(id string, config *configutils.Config, path string, structure any) {
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
	tmpConf.ApplyTransformFuncs(id, structure, "", transformFuncs)

	// Lowercase actual config
	LowercaseKeys(config)

	// Load temporary config back into paths
	config.Layer.Delete(path)
	
	config.Load(tmpConf.Layer.Raw(), path)
}

func InitReload() {
	reload := func(path string) {
		log.Debug(path, " changed, reloading...")
		Load()
		Log()
	}
	
	defaultsConf.OnReload(reload)
	userConf.OnReload(reload)
	tokenConf.OnReload(reload)
}

func InitConfig() {
	var config structure.CONFIG

	mainConf.Layer.Unmarshal("", &config)

	ENV.CONFIGS["*"] = &config

	DEFAULT = ENV.CONFIGS["*"]
}

func LoadDefaults() {
	log.Debug("Loading defaults ", ENV.DEFAULTS_PATH)
	_, err := defaultsConf.LoadFile(ENV.DEFAULTS_PATH, yaml.Parser())

	if err != nil {
		log.Warn("Could not Load Defaults", ENV.DEFAULTS_PATH)
	}
}

func LoadConfig() {
	log.Debug("Loading Config ", ENV.CONFIG_PATH)
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

func normalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", "")

	return key, stringutils.ToType(value)
}
