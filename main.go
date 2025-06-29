package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"

	"main/internal/telegram"
)

type Config struct {
	TelegramToken     string `yaml:"telegram_token"`
	LLMApiKey         string `yaml:"llm_api_key"`
	ObsidianVaultPath string `yaml:"obsidian_vault_path"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	// Load configuration from config.yaml
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	bot, err := telegram.NewBot(cfg.TelegramToken, cfg.LLMApiKey, cfg.ObsidianVaultPath)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}
	bot.Start()
}
