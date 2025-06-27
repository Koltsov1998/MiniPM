package telegram

import (
	"log"
	"os"
	"workspaces/MiniPM/internal/llm"
	"workspaces/MiniPM/internal/obsidian"
)

type Bot struct {
	Token     string
	LLMApiKey string
	VaultPath string
}

func NewBot(token, llmApiKey, vaultPath string) (*Bot, error) {
	return &Bot{
		Token:     token,
		LLMApiKey: llmApiKey,
		VaultPath: vaultPath,
	}, nil
}

func (b *Bot) Start() {
	log.Println("Бот запущен. (Здесь будет логика Telegram API)")
	// TODO: Реализовать обработку сообщений Telegram и интеграцию с LLM и Obsidian
	_ = llm.NewClient(b.LLMApiKey)
	_ = obsidian.NewVault(b.VaultPath)
	os.Exit(0)
}
