package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/codeshelldev/secured-signal-api/utils/safestrings"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var defaultsLayer = koanf.New(".")
var userLayer = koanf.New(".")
var tokensLayer = koanf.New(".")

var config *koanf.Koanf

var configLock sync.Mutex

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

		log.Info(path, " changed, Reloading...")

		configLock.Lock()
		defer configLock.Unlock()

		Load()
	})

	return f, err
}

func LoadDir(path string, dir string, config *koanf.Koanf, parser koanf.Parser) error {
    files, err := filepath.Glob(filepath.Join(dir, "*.yml"))

    if err != nil {
        return err
    }

    for i, file := range files {
		tmp := koanf.New(".")

        _, err := LoadFile(file, tmp, parser)

		if err != nil {
			return err
		}

		config.Set(path + "." + strconv.Itoa(i), tmp.All())
    }

    return nil
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

func mergeConfig(path string, mergeInto *koanf.Koanf, mergeFrom *koanf.Koanf) {
	mergeInto.MergeAt(mergeFrom, path)
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

	config.Delete("")
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
	
	config.Load(confmap.Provider(map[string]any{
		prefix: map[string]any{},
	}, "."), nil)

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