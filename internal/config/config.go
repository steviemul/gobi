package config

import (
	"errors"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

const DEFAULT_PORT = "8080"

type GobiConfig struct {
	port string
}

// fileConfig mirrors the YAML file layout. Fields are pointers/empty so that
// anything omitted from the file keeps its default value.
type fileConfig struct {
	Port string `yaml:"port"`
}

func (gc *GobiConfig) Port() string {
	return gc.port
}

// BuildConfig returns the config described by the YAML file at configPath.
// An empty path or a file that doesn't exist yields the default config; any
// other read or parse failure is returned as an error.
func BuildConfig(configPath string) (GobiConfig, error) {
	config := DefaultConfig()

	if configPath == "" {
		return config, nil
	}

	data, err := os.ReadFile(configPath)

	if errors.Is(err, fs.ErrNotExist) {
		return config, nil
	}

	if err != nil {
		return config, err
	}

	var fc fileConfig

	if err := yaml.Unmarshal(data, &fc); err != nil {
		return config, err
	}

	if fc.Port != "" {
		config.port = fc.Port
	}

	return config, nil
}

func DefaultConfig() GobiConfig {
	return GobiConfig{
		port: DEFAULT_PORT,
	}
}
