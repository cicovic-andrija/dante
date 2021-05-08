package config

// Defaults
const (
	DefaultKeyFile     = ""
	DefaultKey         = ""
	DefaultValidateKey = false
)

// Config
type Config struct {
	KeyFile     string `json:"key_file"`
	Key         string `json:"-"`
	ValidateKey bool   `json:"validate_key"`
}

func newDefaultConfig() *Config {
	return &Config{
		KeyFile:     DefaultKeyFile,
		Key:         DefaultKey,
		ValidateKey: DefaultValidateKey,
	}
}
