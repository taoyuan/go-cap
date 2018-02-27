package config

import (
	"github.com/spf13/cast"
	"strings"
)

type Map struct {
	settings map[string]interface{}
}

func CreateMap(settings map[string]interface{}) Map {
	return Map{
		settings: settings,
	}
}

// GetBool returns the value associated with the key as a boolean.
func (m *Map) GetBool(key string) bool { return cast.ToBool(m.Get(key)) }

// GetString returns the value associated with the key as a string.
func (m *Map) GetString(key string) string { return cast.ToString(m.Get(key)) }

// GetInt returns the value associated with the key as an int.
func (m *Map) GetInt(key string) int { return cast.ToInt(m.Get(key)) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *Map) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(m.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *Map) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(m.Get(key))
}

//  returns the value associated with the key as a slice of strings.
func (m *Map) GetStringSlice(key string) []string {
	return cast.ToStringSlice(m.Get(key))
}

func (m *Map) Get(key string) interface{} {
	if m == nil {
		panic("config not set (key=" + key + ")")
	}
	key = strings.ToLower(key)
	if v, ok := m.settings[key]; ok {
		return v
	}
	return nil
}

func (m *Map) Set(key string, value interface{}) {
	if m == nil {
		panic("language not set")
	}
	key = strings.ToLower(key)
	m.settings[key] = value
}

// IsSet checks whether the key is set in the language or the related config store.
func (m *Map) IsSet(key string) bool {
	key = strings.ToLower(key)

	key = strings.ToLower(key)
	if _, ok := m.settings[key]; ok {
		return true
	}
	return false
}
