package conf

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cicovic-andrija/dante/util"
)

var path string

func init() {
	flag.StringVar(&path, "conf", "", "Config file path")
}

func validateKey(key string) error {
	if util.IsValidUUIDv4(key) {
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
	defer file.Close() // ignore returned error

	scanner := bufio.NewScanner(file)

	// advance the scanner to the EOL
	if !scanner.Scan() {
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

// Load reads settings from a config file
func Load() (*Config, error) {
	if path == "" {
		return nil, errors.New("config file path not provided")
	}

	if _, statErr := os.Stat(path); statErr != nil && os.IsNotExist(statErr) {
		return nil, fmt.Errorf("config file %q not found: %v", path, statErr)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %q: %v", path, err)
	}

	cfg := &Config{}
	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file %q: %v", path, err)
	}

	if err = cfg.Init(); err != nil {
		return nil, fmt.Errorf("%s: failed to initialize config struct: %v", path, err)
	}

	cfg.path = path
	return cfg, nil
}
