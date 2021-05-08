package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	g_configHome string
)

const (
	ConfigDir  = "rac"
	ConfigFile = "config.json"

	dirFilePerm = 0755 // FIXME: is this the correct value?
)

func configPath() (path string, err error) {
	// On Unix systems: $XDG_CONFIG_HOME or $HOME/.config
	// On Windows: %AppData%
	home, err := os.UserConfigDir()
	if err != nil {
		return
	}

	g_configHome = home
	path = filepath.Join(home, ConfigDir, ConfigFile)
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
	defer file.Close() // ignore errors

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
	err := os.MkdirAll(filepath.Join(g_configHome, ConfigDir), dirFilePerm)
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

func Load() (*Config, error) {
	// determine config file location
	path, err := configPath()
	if err != nil {
		err = fmt.Errorf("failed to find config directory: %v", err)
		return nil, err
	}

	// try to write the config file if it doesn't exist
	if _, statErr := os.Stat(path); statErr != nil && os.IsNotExist(statErr) {
		cfg := newDefaultConfig()
		writeConfig(path, cfg) // silently try to write the config file
		return cfg, nil
	}
	// config file exists

	// attempt to open and decode the config file
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("failed to open config file %s: %v", path, err)
		return nil, err
	}

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		err = fmt.Errorf("failed to decode the config file: %v", err)
		return nil, err
	}

	// read API key if key file was specified
	cfg.Key = ""
	if cfg.KeyFile != "" {
		key, err := readKey(cfg.KeyFile, cfg.ValidateKey)
		if err != nil {
			err = fmt.Errorf("failed to read API key: %v", err)
			return nil, err
		}
		cfg.Key = key
	}

	return cfg, nil
}
