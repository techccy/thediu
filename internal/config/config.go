package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("CCY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CCY_API_KEY not set")
	}

	return &Config{
		APIKey:  apiKey,
		BaseURL: getBaseURL(),
		Model:   getModel(),
	}, nil
}

func getBaseURL() string {
	if url := os.Getenv("CCY_API_BASE"); url != "" {
		return url
	}
	return "https://api.openai.com/v1"
}

func getModel() string {
	if model := os.Getenv("CCY_MODEL"); model != "" {
		return model
	}
	return "gpt-4"
}
