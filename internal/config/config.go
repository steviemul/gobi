package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultHTTPPort            = "8080"
	DefaultTLSPort             = "8443"
	LoopbackAddress            = "127.0.0.1"
	DefaultIdleTimeoutSeconds  = 60
	DefaultReadTimeoutSeconds  = 10
	DefaultWriteTimeoutSeconds = 10
)

type GobiConfig struct {
	listen       string
	port         string
	status       bool
	tls          TLS
	idleTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	routes       []Route
}

type TLS struct {
	Enabled bool   `yaml:"enabled"`
	Port    string `yaml:"port"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"key"`
	Force   bool   `yaml:"force"`
}

// Route describes how to handle a request. Host and Path are the match
// criteria (either or both may be set); Root and Proxy are the actions.
// Routes are evaluated in the order they appear in the config file.
type Route struct {
	Host  string `yaml:"host"`
	Path  string `yaml:"path"`
	Index *bool  `yaml:"index"`
	Root  string `yaml:"root"`
	Proxy string `yaml:"proxy"`
}

// fileConfig mirrors the YAML file layout. Anything omitted from the file
// keeps its default value.
type fileConfig struct {
	Listen              string  `yaml:"listen"`
	Port                string  `yaml:"port"`
	Status              *bool   `yaml:"status"`
	TLS                 *TLS    `yaml:"tls"`
	IdleTimeoutSeconds  int     `yaml:"idleTimeoutSeconds"`
	ReadTimeoutSeconds  int     `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int     `yaml:"writeTimeoutSeconds"`
	Routes              []Route `yaml:"routes"`
}

// Listen is the address the HTTP and HTTPS listeners bind to. Empty means
// every interface, on both IPv4 and IPv6.
func (gc *GobiConfig) Listen() string {
	return gc.listen
}

func (gc *GobiConfig) Port() string {
	return gc.port
}

func (gc *GobiConfig) TLS() TLS {
	return gc.tls
}

func (gc *GobiConfig) Status() bool {
	return gc.status
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
		log.Printf("warning: config file %q not found, using defaults", configPath)

		return config, nil
	}

	if err != nil {
		return config, err
	}

	var fc fileConfig

	if err := yaml.Unmarshal(data, &fc); err != nil {
		return config, err
	}

	if fc.Listen != "" {
		config.listen = fc.Listen
	}

	if fc.Port != "" {
		config.port = fc.Port
	}

	if fc.Status != nil {
		config.status = *fc.Status
	}

	if fc.TLS != nil {
		config.tls = *fc.TLS

		if config.tls.Port == "" {
			config.tls.Port = DefaultTLSPort
		}
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
		port:         DefaultHTTPPort,
		status:       true,
		tls:          TLS{Enabled: false, Port: DefaultTLSPort},
		idleTimeout:  DefaultIdleTimeoutSeconds * time.Second,
		readTimeout:  DefaultReadTimeoutSeconds * time.Second,
		writeTimeout: DefaultWriteTimeoutSeconds * time.Second,
		routes:       []Route{},
	}
}

// IndexConfig returns the config behind `gobi -serve`: the default HTTP port
// on the loopback interface only, no TLS, no status page, and a single index route. The route has no path or
// root, so it serves the current working directory at "/".
func IndexConfig() GobiConfig {
	config := DefaultConfig()

	index := true

	config.listen = LoopbackAddress
	config.status = false
	config.routes = []Route{{Index: &index}}

	return config
}
