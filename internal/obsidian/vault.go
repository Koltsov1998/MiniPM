package obsidian

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Vault struct {
	Path string
}

func NewVault(path string) *Vault {
	return &Vault{Path: path}
}

// CreateDailyNote создает новую заметку в хранилище Obsidian с планами на день
func (v *Vault) CreateDailyNote(content string, username string) (string, error) {
	// Проверяем существование директории
	if _, err := os.Stat(v.Path); os.IsNotExist(err) {
		return "", fmt.Errorf("путь к хранилищу Obsidian не существует: %w", err)
	}

	// Создаем директорию для дневниковых записей, если она не существует
	diaryDir := filepath.Join(v.Path, "Дневник")
	if err := os.MkdirAll(diaryDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию для дневника: %w", err)
	}

	// Формируем имя файла на основе текущей даты
	now := time.Now()
	fileName := fmt.Sprintf("%s - %04d-%02d-%02d.md", username, now.Year(), now.Month(), now.Day())
	filePath := filepath.Join(diaryDir, fileName)

	// Формируем содержимое файла
	fileContent := fmt.Sprintf("# Планы на %02d.%02d.%04d\n\n%s\n", 
		now.Day(), now.Month(), now.Year(), content)

	// Записываем файл
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return "", fmt.Errorf("ошибка при записи файла: %w", err)
	}

	return filePath, nil
}

// AppendToDailyNote добавляет содержимое к существующей заметке
func (v *Vault) AppendToDailyNote(filePath, content string) error {
	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("файл не существует: %w", err)
	}

	// Открываем файл для добавления
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	// Добавляем новое содержимое
	_, err = file.WriteString("\n\n" + content)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении содержимого: %w", err)
	}

	return nil
}
