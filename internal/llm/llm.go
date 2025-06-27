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
		BaseURL: "https://api.openai.com/v1/chat/completions",
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// ProcessDailyPlan обрабатывает сообщение пользователя с планами на день
// и возвращает структурированный текст для записи в Obsidian
func (c *Client) ProcessDailyPlan(userMessage string) (string, error) {
	// Формируем запрос к LLM
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
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

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса к LLM: %w", err)
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка API LLM (код %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа LLM: %w", err)
	}

	// Извлекаем текст ответа
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("некорректный формат ответа LLM")
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("некорректный формат сообщения в ответе LLM")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("не удалось извлечь текст из ответа LLM")
	}

	return content, nil
}
