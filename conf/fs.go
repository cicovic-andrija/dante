package conf

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	DefaultConfigDir  = "dante"
	DefaultConfigFile = "config.json"

	NewDirPermissions = 0755
)

func configPath() (path string, err error) {
	// On Unix systems: $XDG_CONFIG_HOME or $HOME/.config
	home, err := os.UserConfigDir()
	if err != nil {
		return
	}

	path = filepath.Join(home, DefaultConfigDir, DefaultConfigFile)
	return
}

func validateKey(key string) error {
	uuidv4 := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	if uuidv4.MatchString(key) {
		return nil
	} else {
		return errors.New("provided key is not a valid UUID v4 string")
	}
}

func readKey(path string, validate bool) (key string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close() // FIXME: handle error

	scanner := bufio.NewScanner(file)

	// advance the scanner to the EOL
	if scanner.Scan() == false {
		// scanning "failed"
		err = scanner.Err()
		if err != nil {
			// scanning really failed
			return
		}

		// scanner encountered EOF, assume line is missing an EOL character
	}

	key = scanner.Text()
	if validate {
		err = validateKey(key)
	}

	return
}

func writeConfig(path string, cfg *Config) error {
	err := os.MkdirAll(filepath.Dir(path), NewDirPermissions)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cfg)
}

// Load reads settings from a config file
func Load(path string) (*Config, error) {

	if path == "" {
		var err error
		path, err = configPath()
		if err != nil {
			return nil, fmt.Errorf("failed to find config directory: %v", err)
		}
	}

	if _, statErr := os.Stat(path); statErr != nil && os.IsNotExist(statErr) {
		cfg := NewDefaultConfig()
		cfg.path = path
		writeConfig(path, cfg) // silently try to write new config file
		return cfg, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %q: %v", path, err)
	}

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file %q: %v", path, err)
	}

	err = cfg.Init()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to initialize config struct: %v", path, err)
	}

	cfg.path = path
	return cfg, nil
}
