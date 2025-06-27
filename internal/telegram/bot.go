package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"main/internal/llm"
	"main/internal/obsidian"
)

type Bot struct {
	Token     string
	LLMApiKey string
	VaultPath string
	api       *tgbotapi.BotAPI
	llmClient *llm.Client
	vault     *obsidian.Vault
}

func NewBot(token, llmApiKey, vaultPath string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации Telegram API: %w", err)
	}

	llmClient := llm.NewClient(llmApiKey)
	vault := obsidian.NewVault(vaultPath)

	return &Bot{
		Token:     token,
		LLMApiKey: llmApiKey,
		VaultPath: vaultPath,
		api:       api,
		llmClient: llmClient,
		vault:     vault,
	}, nil
}

// handleMessage обрабатывает входящее сообщение от пользователя
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Проверяем, что сообщение не пустое
	if strings.TrimSpace(message.Text) == "" {
		log.Println("Получено пустое сообщение, пропускаем")
		return
	}

	// Отправляем сообщение о начале обработки
	processingMsg := tgbotapi.NewMessage(message.Chat.ID, "Обрабатываю ваше сообщение...")
	processingMsg.ReplyToMessageID = message.MessageID
	b.api.Send(processingMsg)

	// Обрабатываем сообщение с помощью LLM
	structuredContent, err := b.llmClient.ProcessDailyPlan(message.Text)
	if err != nil {
		log.Printf("Ошибка обработки сообщения через LLM: %v", err)
		errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при обработке вашего сообщения. Пожалуйста, попробуйте позже.")
		b.api.Send(errorMsg)
		return
	}

	// Сохраняем в Obsidian
	filePath, err := b.vault.CreateDailyNote(structuredContent, message.From.UserName)
	if err != nil {
		log.Printf("Ошибка сохранения в Obsidian: %v", err)
		errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при сохранении в дневник. Пожалуйста, проверьте настройки Obsidian.")
		b.api.Send(errorMsg)
		return
	}

	// Отправляем успешный ответ
	successMsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ваши планы на день успешно сохранены в дневнике!\nФайл: %s", filePath))
	b.api.Send(successMsg)
}

func (b *Bot) Start() {
	log.Printf("Бот запущен: @%s", b.api.Self.UserName)

	// Настраиваем получение обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Обрабатываем сообщения
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Логируем сообщение
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Обрабатываем сообщение
		go b.handleMessage(update.Message)
	}
}
