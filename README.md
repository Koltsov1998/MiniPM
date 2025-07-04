# MiniPM

## Описание
MiniPM - это Telegram-бот, который помогает вести личный дневник в Obsidian. Бот общается с пользователем, обрабатывает сообщения с помощью языковой модели (LLM) и сохраняет записи в формате Markdown в вашем Obsidian vault.

## Возможности
- 💬 Взаимодействие через Telegram: удобный доступ к дневнику через мессенджер
- 🧠 Обработка естественного языка: использование LLM для анализа и структурирования записей
- 📝 Автоматическое создание заметок в Obsidian
- 🔄 Синхронизация между Telegram и Obsidian vault

## Требования
- Go 1.16 или выше
- Telegram Bot Token (получить можно у [@BotFather](https://t.me/BotFather))
- API ключ для OpenRouter (https://openrouter.ai) - сервиса доступа к языковым моделям
- Установленный Obsidian и созданное хранилище (vault)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/MiniPM.git
cd MiniPM
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте конфигурацию (отредактируйте файл `config.yaml`):
```yaml
telegram_token: "YOUR_TELEGRAM_BOT_TOKEN"
# API ключ для OpenRouter (https://openrouter.ai)
llm_api_key: "YOUR_OPENROUTER_API_KEY"
obsidian_vault_path: "/path/to/your/obsidian/vault"
```

4. Скомпилируйте и запустите приложение:
```bash
go build
./MiniPM
```

## Использование

1. Запустите бота
2. Найдите своего бота в Telegram и начните диалог
3. Отправляйте сообщения, которые хотите сохранить в дневнике
4. Бот обработает ваши сообщения и создаст соответствующие записи в вашем Obsidian vault

## Архитектура

Проект состоит из трех основных компонентов:

1. **Telegram Bot** (`internal/telegram`) - обрабатывает взаимодействие с пользователем через Telegram API
2. **LLM Client** (`internal/llm`) - взаимодействует с языковыми моделями через OpenRouter для обработки сообщений
3. **Obsidian Integration** (`internal/obsidian`) - управляет созданием и редактированием Markdown-файлов в Obsidian vault

## Статус проекта

Проект находится в стадии активной разработки. Основные компоненты определены, но требуется реализация конкретной функциональности.

## Лицензия

[MIT License](LICENSE)
