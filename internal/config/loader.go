// internal/config/loader.go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig читает config.yaml и применяет переменные окружения.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Применяем env vars
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.ServerPort = v
	}
	if v := os.Getenv("AI_OCR_URL"); v != "" {
		cfg.AI.OCRURL = v
	}
	if v := os.Getenv("AI_GIGACHAT_URL"); v != "" {
		cfg.AI.GigaChatURL = v
	}
	if v := os.Getenv("AI_GIGACHAT_TOKEN"); v != "" {
		cfg.AI.GigaChatToken = v
	}

	return &cfg, nil
}
