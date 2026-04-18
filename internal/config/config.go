package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProviderConfig struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

type Config struct {
	DefaultProvider string                    `yaml:"default_provider"`
	Providers       map[string]ProviderConfig `yaml:"providers"`
}

var (
	configPath    = filepath.Join(os.Getenv("HOME"), ".ccy", "config.yaml")
	defaultConfig = Config{
		DefaultProvider: "ollama",
		Providers: map[string]ProviderConfig{
			"openai": {
				BaseURL: "https://api.openai.com/v1",
				APIKey:  "env:CCY_OPENAI_KEY",
				Model:   "gpt-4o-mini",
			},
			"deepseek": {
				BaseURL: "https://api.deepseek.com",
				APIKey:  "env:CCY_DEEPSEEK_KEY",
				Model:   "deepseek-chat",
			},
			"ollama": {
				BaseURL: "http://localhost:11434",
				APIKey:  "",
				Model:   "qwen2.5:7b",
			},
		},
	}
)

func Init() error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
	}
	return nil
}

func Load() (*Config, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetCurrentProvider() (string, *ProviderConfig, error) {
	providerName := c.DefaultProvider
	provider, exists := c.Providers[providerName]
	if !exists {
		return "", nil, fmt.Errorf("provider %s not found", providerName)
	}

	if strings.HasPrefix(provider.APIKey, "env:") {
		envKey := strings.TrimPrefix(provider.APIKey, "env:")
		provider.APIKey = os.Getenv(envKey)
	}

	return providerName, &provider, nil
}

func (c *Config) SetDefaultProvider(providerName string) error {
	if _, exists := c.Providers[providerName]; !exists {
		return fmt.Errorf("provider %s not found", providerName)
	}
	c.DefaultProvider = providerName

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) GetProviderNames() []string {
	names := make([]string, 0, len(c.Providers))
	for name := range c.Providers {
		names = append(names, name)
	}
	return names
}

func (c *Config) GetDefaultProvider() string {
	return c.DefaultProvider
}
