package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	data *viper.Viper
}

func New() *Config {
	return &Config{viper.New()}
}

func (c *Config) Data() Provider {
	return c.data
}

func (c *Config) AddFile(file string, silent bool) error {
	c.data.SetConfigFile(file)
	e := c.data.MergeInConfig()
	if e != nil {
		if silent {
			return e
		} else {
			return fmt.Errorf("failed to load file: %s error: %s", file, e)
		}
	}
	return nil
}

func (c *Config) AllSettings() map[string]interface{} {
	return c.data.AllSettings()
}

func (c *Config) GetString(key string) string {
	return c.data.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.data.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.data.GetBool(key)
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.data.GetStringMap(key)
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return c.data.GetStringMapString(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.data.GetStringSlice(key)
}

func (c *Config) Get(key string) interface{} {
	return c.data.Get(key)
}

func (c *Config) Set(key string, value interface{}) {
	c.data.Set(key, value)
}

func (c *Config) IsSet(key string) bool {
	return c.data.IsSet(key)
}
