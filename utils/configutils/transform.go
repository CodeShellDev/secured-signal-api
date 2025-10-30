package configutils

import (
	"maps"
	"reflect"
	"strconv"
	"strings"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/knadh/koanf/providers/confmap"
)

type TransformTarget struct {
	Key string
	Transform string
	ChildTransform string
	Value any
}

func GetKeyToTransformMap(value any) map[string]TransformTarget {
	data := map[string]TransformTarget{}

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

		key := field.Tag.Get("koanf")
		if key == "" {
			continue
		}

		lower := strings.ToLower(key)

		transformTag := field.Tag.Get("transform")
		childTransformTag := field.Tag.Get("childtransform")

		log.Dev("Registering ", lower, " with ", transformTag, " and ", childTransformTag)

		data[lower] = TransformTarget{
			Key:               lower,
			Transform:         transformTag,
			ChildTransform: childTransformTag,
			Value:             getValueSafe(fieldValue),
		}

		// Recursively walk nested structs
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {

			sub := GetKeyToTransformMap(fieldValue.Interface())

			for subKey, subValue := range sub {
				fullKey := lower + "." + strings.ToLower(subKey)

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

	_, res := applyTransform("", all, transformTargets, funcs)

	mapRes, ok := res.(map[string]any)

	if !ok {
		return
	}

	config.Layer.Delete("")
	config.Layer.Load(confmap.Provider(mapRes, "."), nil)
}

func applyTransform(key string, value any, transformTargets map[string]TransformTarget, funcs map[string]func(string, any) (string, any)) (string, any) {
	target := transformTargets[key]

	targets := map[string]TransformTarget{}
		
	maps.Copy(transformTargets, targets)

	newKey, _ := applyTransformToAny(key, value, transformTargets, funcs)

	newKeyWithDot := newKey

	if newKey != "" {
		newKeyWithDot = newKey + "."
	}

	switch asserted := value.(type) {
	case map[string]any:
		res := map[string]any{}

		for k, v := range asserted {
			fullKey := newKeyWithDot + k

			childTarget := TransformTarget{
				Key: fullKey,
				Transform: target.ChildTransform,
				ChildTransform: target.ChildTransform,
			}

			targets[fullKey] = childTarget

			childKey, childValue := applyTransform(fullKey, v, targets, funcs)

			res[childKey] = childValue
		}

		return newKey, res
	case []any:
		res := []any{}
		
		for i, child := range asserted {
			fullKey := newKeyWithDot + strconv.Itoa(i)

			childTarget := TransformTarget{
				Key: strconv.Itoa(i),
				Transform: target.ChildTransform,
				ChildTransform: target.ChildTransform,
			}

			targets[fullKey] = childTarget
			
			_, childValue := applyTransform(fullKey, child, targets, funcs)

			res = append(res, childValue)
		}

		return newKey, res
	default:
		return applyTransformToAny(key, asserted, transformTargets, funcs)
	}
}

func applyTransformToAny(key string, value any, transformTargets map[string]TransformTarget, funcs map[string]func(string, any) (string, any)) (string, any) {
	lower := strings.ToLower(key)

	transformTarget, ok := transformTargets[lower]
	if !ok {
		transformTarget.Transform = "default"
	}

	fn, ok := funcs[transformTarget.Transform]
	if !ok {
		fn = funcs["default"]
	}

	keyParts := getKeyParts(key)

	newKey, newValue := fn(keyParts[len(keyParts)-1], value)
	keyParts[len(keyParts)-1] = newKey

	newFullKey := strings.Join(keyParts, ".")

	log.Dev("Applying ", lower, " with ", transformTarget.Transform, " to ", newFullKey)

	return newFullKey, newValue
}

func getKeyParts(fullKey string) []string {
	keyParts := strings.Split(fullKey, ".")

	return keyParts
}