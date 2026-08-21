package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

const ExampleHash = "$2a$10$example_hash_replace_me"

type Config struct {
	PasswordHash      string `yaml:"password_hash"`
	PasswordFile      string `yaml:"password_file"`
	Port              int    `yaml:"port"`
	DataDir           string `yaml:"data_dir"`
	SessionSecret     string `yaml:"session_secret"`
	SessionSecretFile string `yaml:"session_secret_file"`

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

	if cfg.PasswordFile != "" {
		pwBytes, err := os.ReadFile(cfg.PasswordFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read password_file %s: %w", cfg.PasswordFile, err)
		}
		password := strings.TrimSpace(string(pwBytes))
		if password == "" {
			return nil, fmt.Errorf("password_file %s is empty", cfg.PasswordFile)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password from file: %w", err)
		}
		cfg.PasswordHash = string(hash)
	}

	if cfg.PasswordHash == "" {
		return nil, fmt.Errorf("password_hash or password_file is required in config")
	}

	if cfg.SessionSecretFile != "" {
		secBytes, err := os.ReadFile(cfg.SessionSecretFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read session_secret_file %s: %w", cfg.SessionSecretFile, err)
		}
		cfg.SessionSecret = strings.TrimSpace(string(secBytes))
		if cfg.SessionSecret == "" {
			return nil, fmt.Errorf("session_secret_file %s is empty", cfg.SessionSecretFile)
		}
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("session_secret or session_secret_file is required in config")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	cfg.filePath = absPath

	if cfg.PasswordFile == "" {
		if _, err := bcrypt.Cost([]byte(cfg.PasswordHash)); err != nil {
			return nil, fmt.Errorf("password_hash is not a valid bcrypt hash: %w", err)
		}
		if cfg.PasswordHash == ExampleHash {
			fmt.Println("WARNING: password_hash is still the example value. Change it via the app UI or config.yaml.")
		}
	}

	return &cfg, nil
}

func (c *Config) PasswordExternal() bool {
	return c.PasswordFile != ""
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
