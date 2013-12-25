package config

import (
	"github.com/pelletier/go-toml"
	"path"
)

var tomlconfg *toml.TomlTree
var basePath = "resources/configs/"

// LoadConfig takes an environment name and loads that environment's config file.
func LoadConfig(env string) error {
	var err error
	tomlconfg, err = toml.LoadFile(path.Join(basePath, env+".toml"))
	return err
}

// String returns a config key as a string
func String(key string) string {
	return tomlconfg.Get(key).(string)
}

// Int returns a config key an int64
func Int(key string) int64 {
	return tomlconfg.Get(key).(int64)
}
