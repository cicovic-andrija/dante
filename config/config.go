package config

import "fmt"

// Defaults
const (
	DefaultKeyFile     = ""
	DefaultValidateKey = false
)

// Config
type Config struct {
	Auth Auth `json:"auth"`
}

// Auth
type Auth struct {
	KeyFile     string `json:"key_file"`
	Key         string `json:"-"`
	ValidateKey bool   `json:"validate_key"`
}

func NewDefaultConfig() *Config {
	cfg := &Config{
		Auth: Auth{
			KeyFile:     DefaultKeyFile,
			ValidateKey: DefaultValidateKey,
		},
	}
	cfg.Init() // must not fail here
	return cfg
}

func (cfg *Config) Init() error {

	cfg.Auth.Key = ""
	if cfg.Auth.KeyFile != "" {
		key, err := readKey(cfg.Auth.KeyFile, cfg.Auth.ValidateKey)
		if err != nil {
			return fmt.Errorf("failed to read API key: %v", err)
		}
		cfg.Auth.Key = key
	}

	return nil
}
