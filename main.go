package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

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

// loadEnvLocal loads configuration from env.local file
func loadEnvLocal() map[string]string {
	env := make(map[string]string)

	file, err := os.Open("env.local")
	if err != nil {
		log.Printf("env.local не найден, используем только config.yaml")
		return env
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			env[key] = value
		}
	}

	return env
}

func main() {
	// Load configuration from config.yaml
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Override with env.local if available
	env := loadEnvLocal()
	if telegramToken, ok := env["telegram_token"]; ok {
		cfg.TelegramToken = telegramToken
		log.Println("Используем токен Telegram из env.local")
	}

	bot, err := telegram.NewBot(cfg.TelegramToken, cfg.LLMApiKey, cfg.ObsidianVaultPath)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}
	bot.Start()
}
