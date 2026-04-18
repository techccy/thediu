package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ccy-ai/ccy-assistant/internal/config"
)

type Client struct {
	apiKey  string
	baseURL string
	model   string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
	}
}

func (c *Client) SendRequest(failedCommand, errorMessage string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("CCY_API_KEY environment variable is not set")
	}

	systemPrompt := `You are a Linux/Unix system expert. Your task is to analyze the failed command and error message, then provide a corrected command.

You must respond with a valid JSON object only, without any markdown code blocks or extra text. The JSON must have this exact format:
{
  "analysis": "Brief explanation of the error cause in one sentence",
  "command": "The corrected command that should be executed"
}

Example:
{
  "analysis": "The command 'pushu' is not a valid git subcommand",
  "command": "git push origin main"
}`

	var userContent string
	if errorMessage != "" {
		userContent = fmt.Sprintf("Failed command: %s\nError message: %s", failedCommand, errorMessage)
	} else {
		userContent = fmt.Sprintf("Failed command: %s", failedCommand)
	}

	reqBody := Request{
		Model: c.model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
