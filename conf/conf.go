package conf

import (
	"fmt"
	"os"
)

// Defaults
const (
	DefaultEnv         = "dev"
	DefaultTcpPort     = 8080
	DefaultKeyFile     = ""
	DefaultValidateKey = false
	DefaultLogDir      = "/tmp/log"
)

// Limits
const (
	MaxPortNumber = 65535
)

// Config
type Config struct {
	Env  string `json:"env"`
	Net  Net    `json:"net"`
	Auth Auth   `json:"auth"`
	Log  Log    `json:"log"`

	path string `json:"-"`
}

// Net
type Net struct {
	Tcp Tcp `json:"tcp"`
}

// Tcp
type Tcp struct {
	Port int `json:"port"`
}

// Auth
type Auth struct {
	KeyFile     string `json:"key_file"`
	Key         string `json:"-"`
	ValidateKey bool   `json:"validate_key"`
}

// Log
type Log struct {
	Dir string `json:"dir"`
}

func NewDefaultConfig() *Config {
	cfg := &Config{
		Env: DefaultEnv,
		Net: Net{
			Tcp: Tcp{
				Port: DefaultTcpPort,
			},
		},
		Auth: Auth{
			KeyFile:     DefaultKeyFile,
			ValidateKey: DefaultValidateKey,
		},
		Log: Log{
			Dir: DefaultLogDir,
		},
	}
	cfg.Init() // must not fail here
	return cfg
}

func (cfg *Config) Path() string {
	return cfg.path
}

func (cfg *Config) Init() error {

	if cfg.Env == "" {
		return fmt.Errorf("invalid environment: env cannot be empty")
	}

	if cfg.Net.Tcp.Port < 1 || cfg.Net.Tcp.Port > MaxPortNumber {
		return fmt.Errorf("invalid TCP port number")
	}

	cfg.Auth.Key = ""
	if cfg.Auth.KeyFile != "" {
		key, err := readKey(cfg.Auth.KeyFile, cfg.Auth.ValidateKey)
		if err != nil {
			return fmt.Errorf("failed to read API key: %v", err)
		}
		cfg.Auth.Key = key
	}

	if finfo, statErr := os.Stat(cfg.Log.Dir); statErr != nil && os.IsNotExist(statErr) {
		return fmt.Errorf("log path %q doesn't exist", cfg.Log.Dir)
	} else if !finfo.IsDir() {
		return fmt.Errorf("log path %q is not a directory", cfg.Log.Dir)
	}

	return nil
}
