package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	ApiKey     string
	BaseURL    string
	HttpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey:  apiKey,
		BaseURL: "https://openrouter.ai/api/v1/chat/completions",
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// ProcessDailyPlan обрабатывает сообщение пользователя с планами на день
// и возвращает структурированный текст для записи в Obsidian
func (c *Client) ProcessDailyPlan(userMessage string) (string, error) {
	// Формируем запрос к LLM через OpenRouter
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "openai/gpt-3.5-turbo", // Указываем провайдера и модель
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты помощник, который структурирует планы на день для записи в дневник. Форматируй текст в Markdown с заголовками, списками и выделением важных моментов.",
			},
			{
				"role":    "user",
				"content": userMessage,
			},
		},
		"temperature": 0.7,
		"route": "fallback", // Использовать запасную модель, если основная недоступна
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при формировании запроса к LLM: %w", err)
	}

	// Отправляем запрос
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании HTTP запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("HTTP-Referer", "https://minipm.app") // Идентификатор приложения для OpenRouter
	req.Header.Set("X-Title", "MiniPM") // Название приложения для OpenRouter

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса к LLM: %w", err)
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка API OpenRouter (код %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Читаем тело ответа для логирования и повторного использования
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа OpenRouter: %w", err)
	}

	// Декодируем JSON ответ
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа OpenRouter: %w", err)
	}

	// Проверяем наличие ошибки в ответе
	if errMsg, hasError := result["error"].(map[string]interface{}); hasError {
		errMessage := "неизвестная ошибка"
		if msg, ok := errMsg["message"].(string); ok {
			errMessage = msg
		}
		return "", fmt.Errorf("ошибка OpenRouter: %s", errMessage)
	}

	// Извлекаем текст ответа
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("некорректный формат ответа OpenRouter: отсутствует поле 'choices' или оно пустое")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("некорректный формат первого выбора в ответе OpenRouter")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("некорректный формат сообщения в ответе OpenRouter")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("не удалось извлечь текст из ответа OpenRouter")
	}

	return content, nil
}
