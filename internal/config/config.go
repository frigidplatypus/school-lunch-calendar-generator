package config

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

const ExampleHash = "$2a$10$example_hash_replace_me"

type Config struct {
	PasswordHash  string `yaml:"password_hash"`
	Port          int    `yaml:"port"`
	DataDir       string `yaml:"data_dir"`
	SessionSecret string `yaml:"session_secret"`

	filePath string
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	if cfg.PasswordHash == "" {
		return nil, fmt.Errorf("password_hash is required in config")
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("session_secret is required in config")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	cfg.filePath = absPath

	if _, err := bcrypt.Cost([]byte(cfg.PasswordHash)); err != nil {
		return nil, fmt.Errorf("password_hash is not a valid bcrypt hash: %w", err)
	}

	if cfg.PasswordHash == ExampleHash {
		fmt.Println("WARNING: password_hash is still the example value. Change it via the app UI or config.yaml.")
	}

	return &cfg, nil
}

func (c *Config) UpdatePasswordHash(newHash string) error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot parse config file: %w", err)
	}

	raw["password_hash"] = newHash
	c.PasswordHash = newHash

	out, err := yaml.Marshal(raw)
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}

	return os.WriteFile(c.filePath, out, 0600)
}

func (c *Config) LogDir() string {
	return filepath.Join(c.DataDir, "logos")
}

func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "calendar-data.db")
}
