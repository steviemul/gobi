package config

import (
	"errors"
	"io/fs"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_PORT                  = "8080"
	DEFAULT_IDLE_TIMEOUT_SECONDS  = 60
	DEFAULT_READ_TIMEOUT_SECONDS  = 10
	DEFAULT_WRITE_TIMEOUT_SECONDS = 10
)

type GobiConfig struct {
	port         string
	status       bool
	idleTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	routes       []Route
}

// Route describes how to handle a request. Host and Path are the match
// criteria (either or both may be set); Root and Proxy are the actions.
// Routes are evaluated in the order they appear in the config file.
type Route struct {
	Host  string `yaml:"host"`
	Path  string `yaml:"path"`
	Root  string `yaml:"root"`
	Proxy string `yaml:"proxy"`
}

// fileConfig mirrors the YAML file layout. Anything omitted from the file
// keeps its default value.
type fileConfig struct {
	Port                string  `yaml:"port"`
	Status              bool    `yaml:"status"`
	IdleTimeoutSeconds  int     `yaml:"idleTimeoutSeconds"`
	ReadTimeoutSeconds  int     `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int     `yaml:"writeTimeoutSeconds"`
	Routes              []Route `yaml:"routes"`
}

func (gc *GobiConfig) Port() string {
	return gc.port
}

// IdleTimeout is the max wait between requests on a keep-alive connection.
func (gc *GobiConfig) IdleTimeout() time.Duration {
	return gc.idleTimeout
}

// ReadTimeout is the max time to receive one request.
func (gc *GobiConfig) ReadTimeout() time.Duration {
	return gc.readTimeout
}

// WriteTimeout is the max time to send one response.
func (gc *GobiConfig) WriteTimeout() time.Duration {
	return gc.writeTimeout
}

func (gc *GobiConfig) Routes() []Route {
	return gc.routes
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

	if fc.IdleTimeoutSeconds > 0 {
		config.idleTimeout = time.Duration(fc.IdleTimeoutSeconds) * time.Second
	}

	if fc.ReadTimeoutSeconds > 0 {
		config.readTimeout = time.Duration(fc.ReadTimeoutSeconds) * time.Second
	}

	if fc.WriteTimeoutSeconds > 0 {
		config.writeTimeout = time.Duration(fc.WriteTimeoutSeconds) * time.Second
	}

	if fc.Routes != nil {
		config.routes = fc.Routes
	}

	return config, nil
}

func DefaultConfig() GobiConfig {
	return GobiConfig{
		port:         DEFAULT_PORT,
		status:       true,
		idleTimeout:  DEFAULT_IDLE_TIMEOUT_SECONDS * time.Second,
		readTimeout:  DEFAULT_READ_TIMEOUT_SECONDS * time.Second,
		writeTimeout: DEFAULT_WRITE_TIMEOUT_SECONDS * time.Second,
		routes:       []Route{},
	}
}
