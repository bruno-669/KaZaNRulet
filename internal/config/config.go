// internal/config/config.go
package config

type Config struct {
	ServerPort string   `yaml:"server_port"`
	AI         AIConfig `yaml:"ai"`
}

type AIConfig struct {
	OCRURL        string `yaml:"ocr_url"`
	GigaChatURL   string `yaml:"gigachat_url"`
	GigaChatToken string `yaml:"gigachat_token"`
}
