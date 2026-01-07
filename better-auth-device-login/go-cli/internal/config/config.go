package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

const (
	configDirName  = "go-auth-device-cli"
	configFileName = "config.toml"
)

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	Hostname  string    `toml:"hostname"`
	UserEmail string    `toml:"user_email,omitempty"`
	ExpiresAt time.Time `toml:"expires_at,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Auth AuthConfig `toml:"auth"`
}

// configDir returns the path to the config directory
func configDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", configDirName), nil
}

// configPath returns the path to the config file
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0700)
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func Save(cfg *Config) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// Clear removes the configuration file
func Clear() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
