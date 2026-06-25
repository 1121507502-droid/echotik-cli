package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultBaseURL = "https://open.echotik.live"
	envConfigDir   = "ECHOTIK_CLI_CONFIG_DIR"
)

type Config struct {
	BaseURL  string `json:"baseUrl"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func ConfigDir() string {
	if dir := strings.TrimSpace(os.Getenv(envConfigDir)); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".echotik-cli"
	}
	return filepath.Join(home, ".echotik-cli")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if err := os.MkdirAll(ConfigDir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), append(data, '\n'), 0600)
}

func ResolveCredential() (baseURL, username, password string, err error) {
	baseURL = strings.TrimSpace(os.Getenv("ECHOTIK_BASE_URL"))
	username = strings.TrimSpace(os.Getenv("ECHOTIK_USERNAME"))
	password = strings.TrimSpace(os.Getenv("ECHOTIK_PASSWORD"))

	cfg, cfgErr := Load()
	if cfgErr == nil {
		if baseURL == "" {
			baseURL = cfg.BaseURL
		}
		if username == "" {
			username = cfg.Username
		}
		if password == "" {
			password = cfg.Password
		}
	} else if !errors.Is(cfgErr, os.ErrNotExist) {
		return "", "", "", cfgErr
	}

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if username == "" || password == "" {
		return baseURL, username, password, errors.New("missing EchoTik username or password")
	}
	return baseURL, username, password, nil
}
