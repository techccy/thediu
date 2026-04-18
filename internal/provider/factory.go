package provider

import (
	"fmt"
	"os"

	"github.com/techccy/diu-assistant/internal/config"
)

type ProviderFactory struct{}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

func (f *ProviderFactory) CreateProvider(cfg *config.Config, providerName string) (LLMProvider, error) {
	provider, exists := cfg.Providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	apiKey := provider.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("DIU_API_KEY")
	}

	switch providerName {
	case "openai":
		return NewOpenAIProvider(apiKey, provider.BaseURL, provider.Model), nil
	case "deepseek":
		return NewDeepSeekProvider(apiKey, provider.BaseURL, provider.Model), nil
	case "ollama":
		return NewOllamaProvider(provider.BaseURL, provider.Model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}
