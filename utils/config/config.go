package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codeshelldev/secured-signal-api/utils"
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

	WatchFile(path, f)

	return f, err
}

func WatchFile(path string, f *file.File) {
	f.Watch(func(event any, err error) {
		if err != nil {
			return
		}

		log.Info(path, " changed, Reloading...")

		configLock.Lock()
		defer configLock.Unlock()

		Load()
	})
}

func LoadDir(path string, dir string, config *koanf.Koanf, parser koanf.Parser) error {
    files, err := filepath.Glob(filepath.Join(dir, "*.yml"))

    if err != nil {
        return nil
    }

	var array []any

	for _, f := range files {
		tmp := koanf.New(".")

		_, err := LoadFile(f, tmp, parser)

		if err != nil {
			return err
		}

		array = append(array, tmp.Raw())
	}

	wrapper := map[string]any{
		path: array,
	}

    return config.Load(confmap.Provider(wrapper, "."), nil)
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

// Transforms Children of path
func transformChildren(config *koanf.Koanf, path string, transform func(key string, value any) (string, any)) error {
	var sub map[string]any
	
	if !config.Exists(path) {
		return errors.New("invalid path")
	}

	err := config.Unmarshal(path, &sub)
	
	if err != nil {
		return err
	}

	transformed := make(map[string]any)

	for key, val := range sub {
		newKey, newVal := transform(key, val)

		transformed[newKey] = newVal
	}
	
	config.Delete(path)

	config.Load(confmap.Provider(map[string]any{
		path: transformed,
	}, "."), nil)

	return nil
}

// Does the same thing as transformChildren() but does it for each Array Item inside of root and transforms subPath
func transformChildrenUnderArray(config *koanf.Koanf, root string, subPath string, transform func(key string, value any) (string, any)) error {
	var array []map[string]any
	
	err := config.Unmarshal(root, &array)
	if err != nil {
		return err
	}

	transformed := []map[string]any{}

	for _, data := range array {
		tmp := koanf.New(".")

		tmp.Load(confmap.Provider(map[string]any{
			"item": data,
		}, "."), nil)

		log.Dev(utils.ToJson(tmp.All()))

		err := transformChildren(tmp, "item." + subPath, transform)

		if err != nil {
			return err
		}

		log.Dev(utils.ToJson(tmp.All()))

		item := tmp.Get("item")

		if item != nil {
			itemMap, ok := item.(map[string]any)

			if ok {
				transformed = append(transformed, itemMap)
			}
		}
	}

	config.Load(confmap.Provider(map[string]any{
		root: map[string]any{},
	}, "."), nil)

	config.Load(confmap.Provider(map[string]any{
		root: transformed,
	}, "."), nil)

	return nil
}


func normalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", "")

	return key, safestrings.ToType(value)
}