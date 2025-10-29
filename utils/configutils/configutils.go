package configutils

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
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

type TransformTarget struct {
	Key string
	Transform string
	Value any
}

func New() *Config {
	return &Config{
		Layer: koanf.New("."),
		LoadFunc: func() {
			log.Dev("Config.LoadFunc not initialized!")
		},
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

func GetKeyToTransformMap(value any) map[string]TransformTarget {
	data := map[string]TransformTarget{}

	log.Info("Value: ", jsonutils.ToJson(value))

	if value == nil {
		return data
	}

	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return data
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		log.Info("Field: ", field.Name)

		key := field.Tag.Get("koanf")
		if key == "" {
			continue
		}

		transformTag := field.Tag.Get("transform")

		data[key] = TransformTarget{
			Key:       key,
			Transform: transformTag,
			Value:     getValueSafe(fieldValue),
		}

		log.Info(key, ": ", v.String())

		// Recursively walk nested structs
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {

			sub := GetKeyToTransformMap(fieldValue.Interface())

			for subKey, subValue := range sub {
				fullKey := key + "." + subKey

				data[fullKey] = subValue
			}
		}
	}

	return data
}

func getValueSafe(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return getValueSafe(value.Elem())
	}
	return value.Interface()
}

func (config Config) ApplyTransformFuncs(structSchema any, path string, funcs map[string]func(string, any) (string, any)) {
	transformTargets := GetKeyToTransformMap(structSchema)

	var all map[string]any

	if path == "." {
		all = config.Layer.All()
	} else {
		all = config.Layer.Get(path).(map[string]any)
	}

	res := map[string]any{}

	for key, value := range all {
		transformTarget, ok := transformTargets[key]
		if !ok {
			transformTarget.Transform = "default"
		}

		fn, ok := funcs[transformTarget.Transform]
		if !ok {
			fn = funcs["default"]
		}

		log.Info()

		newKey, newValue := fn(key, value)

		res[newKey] = newValue
	}

	config.Layer.Delete("")
	config.Layer.Load(confmap.Provider(res, "."), nil)
}

func WatchFile(path string, f *file.File, loadFunc func()) {
	f.Watch(func(event any, err error) {
		if err != nil {
			return
		}

		log.Info(path, " changed, Reloading...")

		configLock.Lock()
		defer configLock.Unlock()

		loadFunc()
	})
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

// Transforms Children of path
func (config *Config) TransformChildren(path string, transform func(key string, value any) (string, any)) error {
	var sub map[string]any

	if !config.Layer.Exists(path) {
		return errors.New("invalid path")
	}

	err := config.Layer.Unmarshal(path, &sub)

	if err != nil {
		return err
	}

	transformed := make(map[string]any)

	for key, val := range sub {
		newKey, newVal := transform(key, val)

		transformed[newKey] = newVal
	}

	config.Layer.Delete(path)

	config.Layer.Load(confmap.Provider(map[string]any{
		path: transformed,
	}, "."), nil)

	return nil
}

// Does the same thing as transformChildren() but does it for each Array Item inside of root and transforms subPath
func (config *Config) TransformChildrenUnderArray(root string, subPath string, transform func(key string, value any) (string, any)) error {
	var array []map[string]any

	err := config.Layer.Unmarshal(root, &array)
	if err != nil {
		return err
	}

	transformed := []map[string]any{}

	for _, data := range array {
		tmp := New()

		tmp.Layer.Load(confmap.Provider(map[string]any{
			"item": data,
		}, "."), nil)

		err := tmp.TransformChildren("item."+subPath, transform)

		if err != nil {
			return err
		}

		item := tmp.Layer.Get("item")

		if item != nil {
			itemMap, ok := item.(map[string]any)

			if ok {
				transformed = append(transformed, itemMap)
			}
		}
	}

	config.Layer.Delete(root)

	config.Layer.Load(confmap.Provider(map[string]any{
		root: transformed,
	}, "."), nil)

	return nil
}

func (config *Config) NormalizeEnv(key string, value string) (string, any) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", "")

	return key, stringutils.ToType(value)
}
