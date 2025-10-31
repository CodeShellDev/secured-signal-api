package configutils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	stringutils "github.com/codeshelldev/secured-signal-api/utils/stringutils"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var configLock sync.Mutex

type Config struct {
	Layer *koanf.Koanf
	LoadFunc func()
}

func New() *Config {
	return &Config{
		Layer: koanf.New("."),
		LoadFunc: func() {},
	}
}

func (config *Config) OnLoad(onLoad func()) {
	config.LoadFunc = onLoad
}

func (config *Config) LoadFile(path string, parser koanf.Parser) (koanf.Provider, error) {
	log.Debug("Loading Config File: ", path)

	f := file.Provider(path)

	err := config.Layer.Load(f, parser)
	
	if err != nil {
		return nil, err
	}

	WatchFile(path, f, config.LoadFunc)

	return f, err
}

func WatchFile(path string, f *file.File, loadFunc func()) {
	f.Watch(func(event any, err error) {
		if err != nil {
			return
		}

		log.Info(path, " changed, Reloading...")

		configLock.Lock()
		defer configLock.Unlock()

		f.Unwatch()

		loadFunc()
	})
}

func getPath(str string) string {
	if str == "." {
		str = ""
	}

	return str
}

func (config *Config) Delete(path string) (error) {
	if !config.Layer.Exists(path) {
		return errors.New("path not found")
	}

	all := config.Layer.Get("")

	log.Dev("Init:\n--------------------------------------\n", jsonutils.ToJson(all), "\n--------------------------------------")
	
	if all == nil {
		return errors.New("empty config")
	}

	for _, key := range config.Layer.Keys() {
		if strings.HasPrefix(key, path + ".") || key == path {
			config.Layer.Delete(key)
		}
	}

	log.Dev("Deletion:\n--------------------------------------\n", jsonutils.ToJson(all), "\n--------------------------------------")

	return nil
}

func (config *Config) LoadDir(path string, dir string, ext string, parser koanf.Parser) error {
	files, err := filepath.Glob(filepath.Join(dir, "*" + ext))

	if err != nil {
		return nil
	}

	var array []any

	for _, f := range files {
		tmp := New()

		_, err := tmp.LoadFile(f, parser)

		if err != nil {
			return err
		}

		array = append(array, tmp.Layer.Raw())
	}

	wrapper := map[string]any{
		path: array,
	}

	return config.Layer.Load(confmap.Provider(wrapper, "."), nil)
}

func (config *Config) LoadEnv() (koanf.Provider, error) {
	e := env.Provider(".", env.Opt{
		TransformFunc: config.NormalizeEnv,
	})

	err := config.Layer.Load(e, nil)

	if err != nil {
		log.Fatal("Error loading env: ", err.Error())
	}

	return e, err
}

func (config *Config) TemplateConfig() {
	data := config.Layer.All()

	for key, value := range data {
		str, isStr := value.(string)

		if isStr {
			templated := os.ExpandEnv(str)

			if templated != "" {
				data[key] = templated
			}
		}
	}

	config.Layer.Load(confmap.Provider(data, "."), nil)
}

func (config *Config) MergeLayers(layers ...*koanf.Koanf) {
	for _, layer := range layers {
		config.Layer.Merge(layer)
	}
}

func (config *Config) NormalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", "")

	return key, stringutils.ToType(value)
}
