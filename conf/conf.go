package conf

import (
	"errors"
	"fmt"
	"os"
)

const (
	MaxPortNumber = 65535
	HTTPString    = "http"
	HTTPSString   = "https"
)

// Config
type Config struct {
	Env    string       `json:"env"`
	Net    Net          `json:"net"`
	Atlas  AtlasConf    `json:"atlas"`
	Influx InfluxDBConf `json:"influxdb"`
	Log    Log          `json:"log"`

	path string `json:"-"`
}

// Net
type Net struct {
	Protocol string `json:"protocol"`
	DNSName  string `json:"dns_name"`
	Port     int    `json:"port"`
}

// AtlasConf
type AtlasConf struct {
	Net  Net       `json:"net"`
	Auth AtlasAuth `json:"auth"`
}

// AtlasAuth
type AtlasAuth struct {
	KeyFile     string `json:"key_file"`
	Key         string `json:"-"`
	ValidateKey bool   `json:"validate_key"`
}

// InfluxDB
type InfluxDBConf struct {
	Net Net `json:"net"`
}

// Log
type Log struct {
	Dir string `json:"dir"`
}

func (cfg *Config) Path() string {
	return cfg.path
}

func (cfg *Config) Init() error {
	var err error

	if cfg.Env == "" {
		return fmt.Errorf("invalid environment: env cannot be empty")
	}

	if err = validateNetConf(&cfg.Net); err != nil {
		return err
	}

	if err = validateNetConf(&cfg.Atlas.Net); err != nil {
		return err
	}

	cfg.Atlas.Auth.Key = ""
	if cfg.Atlas.Auth.KeyFile != "" {
		key, err := readKey(cfg.Atlas.Auth.KeyFile, cfg.Atlas.Auth.ValidateKey)
		if err != nil {
			return fmt.Errorf("failed to read Atlas API key: %v", err)
		}
		cfg.Atlas.Auth.Key = key
	} // FIXME: What if key file is not specified?

	if finfo, statErr := os.Stat(cfg.Log.Dir); statErr != nil && os.IsNotExist(statErr) {
		return fmt.Errorf("log path %q doesn't exist", cfg.Log.Dir)
	} else if !finfo.IsDir() {
		return fmt.Errorf("log path %q is not a directory", cfg.Log.Dir)
	}

	return nil
}

func validateNetConf(net *Net) error {
	const errorPrefix = "net config validation failed: "

	if net.Protocol != HTTPString && net.Protocol != HTTPSString {
		return fmt.Errorf("%sinvalid protocol %q", errorPrefix, net.Protocol)
	}

	if net.DNSName == "" {
		return errors.New(errorPrefix + "empty DNS name")
	}

	if net.Port < 1 || net.Port > MaxPortNumber {
		return fmt.Errorf("%sinvalid TCP port number: %d", errorPrefix, net.Port)
	}

	return nil
}
