package config

import (
	"errors"
	"reflect"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	"github.com/codeshelldev/secured-signal-api/utils/prettylog"
)

type TestableWithRaw interface {
	TestWithRaw(path string, raw any) error
}

type Testable interface {
	Test(path string) error
}

func Validate() {
	for _, config := range ENV.CONFIGS {
		configutils.WalkSchema(reflect.TypeFor[structure.CONFIG](), reflect.ValueOf(*config), config.RAW.Layer.Get(""), nil, func(path string, field reflect.StructField, raw any, value reflect.Value) {
			tryTest(path, raw, field.Type, value)
		})
	}
}

func tryTest(path string, raw any, schema reflect.Type, value reflect.Value) {
	ifaceType1 := reflect.TypeFor[TestableWithRaw]()
	ifaceType2 := reflect.TypeFor[Testable]()

	// unwrap pointers
	for schema.Kind() == reflect.Pointer {
		schema = schema.Elem()
	}

	if !schema.Implements(ifaceType1) && !reflect.PointerTo(schema).Implements(ifaceType1) && !schema.Implements(ifaceType2) && !reflect.PointerTo(schema).Implements(ifaceType2) {
		return
	}

	// ensure value exists for execution
	if !value.IsValid() {
		if schema.Kind() == reflect.Pointer {
			value = reflect.New(schema.Elem())
		} else {
			value = reflect.New(schema).Elem()
		}
	}

	if !value.CanInterface() {
		return
	}

	tRaw, ok := value.Interface().(TestableWithRaw)

	if ok {
		err := tRaw.TestWithRaw(path, raw)

		if err != nil {
			prettylog.GenericError("{b,fg=red}Config syntax Error{/}", errors.New("syntax error at '" + path + "':\n\n" + err.Error()))
		}
	}

	t, ok := value.Interface().(Testable)

	if ok {
		err := t.Test(path)

		if err != nil {
			prettylog.GenericError("{b,fg=red}Config syntax Error{/}", errors.New("syntax error at '" + path + "':\n\n" + err.Error()))
		}
	}
}