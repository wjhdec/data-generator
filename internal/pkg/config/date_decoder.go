package config

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

func DateDecoder(m *mapstructure.DecoderConfig) {
	m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		func(
			f reflect.Type,
			t reflect.Type,
			data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(time.Time{}) {
				return data, nil
			}

			asString := data.(string)
			if asString == "" {
				return time.Time{}, nil
			}

			return time.Parse("2006-01-02 15:04:05", asString)
		},
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
}
