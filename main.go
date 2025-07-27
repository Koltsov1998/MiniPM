package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"main/config"
	"os"

	"main/internal/telegram"
)

type MongoDBOptions struct {
	MaxPoolSize      int `yaml:"max_pool_size"`
	ConnectTimeoutMS int `yaml:"connect_timeout_ms"`
}

type MongoDBConfig struct {
	URI      string         `yaml:"uri"`
	Database string         `yaml:"database"`
	Options  MongoDBOptions `yaml:"options"`
}

type Config struct {
	TelegramToken     string        `yaml:"telegram_token"`
	LLMApiKey         string        `yaml:"llm_api_key"`
	ObsidianVaultPath string        `yaml:"obsidian_vault_path"`
	MongoDB           MongoDBConfig `yaml:"mongodb"`
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
	var cfg *Config

	if _, err := os.Stat("config.local.yaml"); err == nil {

		log.Println("Загрузка основной конфигурации из config.yaml")
		baseCfg, err := loadConfig("config.yaml")
		if err != nil {
			log.Fatalf("Ошибка загрузки основной конфигурации: %v", err)
		}
		cfg = baseCfg
	}

	if _, err := os.Stat("config.local.yaml"); err == nil {
		log.Println("Загрузка и объединение локальной конфигурации из config.local.yaml")
		localCfg, err := loadConfig("config.local.yaml")
		if err != nil {
			log.Printf("Ошибка загрузки локальной конфигурации: %v, используем только основную конфигурацию", err)
		} else {
			err = config.MergeObjects(cfg, localCfg)
			if err != nil {
				log.Println("Конфигурации успешно объединены")
			}
		}
	}

	if cfg == nil {
		log.Fatalf("Failed to load config.yaml and config.local.yaml. Please check your configuration and try again.")
	}

	bot, err := telegram.NewBot(cfg.TelegramToken, cfg.LLMApiKey, cfg.ObsidianVaultPath)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}
	bot.Start()
}
