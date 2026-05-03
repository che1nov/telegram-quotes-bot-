# Telegram Quotes Bot

[English version](README.md)

Telegram-бот, который по расписанию отправляет мотивационные цитаты на русском языке в Telegram-чат или канал. Бот получает случайные цитаты из Forismatic API и отправляет их через Telegram Bot API.

## Возможности
- Получение случайных мотивационных цитат на русском языке из Forismatic API
- Отправка цитат в Telegram по cron-расписанию
- Опциональная отправка тестовой цитаты при запуске
- Graceful shutdown по системным сигналам
- HTTP-клиенты с таймаутами
- Валидация конфигурации и данных из внешнего API
- Структурированное логирование через `log/slog`
- Простая слоистая структура с use cases и adapters

## Технологии
- **Язык**: Go 1.23
- **API**: Telegram Bot API, Forismatic API
- **Планировщик**: robfig/cron
- **Архитектура**: простая слоистая структура с dependency injection
- **Тестирование**: стандартный пакет Go testing
- **Контейнеризация**: Docker
- **Task runner**: Taskfile

## Структура проекта
```text
├── cmd/                   # Точка входа приложения
├── internal/
│   ├── adapters/          # Адаптеры внешних сервисов
│   ├── config/            # Загрузка конфигурации
│   ├── entities/          # Доменные сущности
│   ├── usecases/          # Сценарии приложения
│   └── validators/        # Валидация входных данных
├── Dockerfile             # Docker-образ
├── Taskfile.yml           # Автоматизация задач
└── README.md              # README на английском
```

## Быстрый старт

### Требования
- Go 1.23+
- Docker, опционально
- Task, опционально

### Настройка
1. Склонируйте репозиторий.
2. Создайте файл `.env`.
3. Укажите токен Telegram-бота и ID целевого чата:
   ```bash
   BOT_TOKEN=123456789:your-token
   CHAT_ID=-1001234567890
   SEND_TEST_QUOTE=true
   ```

### Запуск
```bash
# Go
go run cmd/main.go

# Task
task run

# Docker
docker build -t telegram-quotes-bot .
docker run --env-file .env telegram-quotes-bot
```

### Разработка
```bash
# Запустить тесты
task test

# Запустить тесты с покрытием
task test-coverage

# Собрать приложение
task build

# Запустить go vet и gofmt
task lint

# Удалить артефакты сборки
task clean
```

## Конфигурация
Бот читает конфигурацию из переменных окружения:
- `BOT_TOKEN`: токен Telegram-бота от BotFather
- `CHAT_ID`: ID целевого Telegram-чата или канала
- `SEND_TEST_QUOTE`: отправлять тестовую цитату при запуске, по умолчанию `true`

## Расписание
Бот отправляет цитаты каждый день в `04:00`, `08:00`, `14:00` и `18:00`.

## API
Бот использует Forismatic API для получения русскоязычных цитат:
- **Endpoint**: `http://api.forismatic.com/api/1.0/?method=getQuote&format=json&lang=ru`
- **Поля ответа**: `quoteText`, `quoteAuthor`
- **Язык**: русский

## Безопасность
- Токен бота маскируется в логах
- HTTP-клиенты используют таймауты
- Внешние данные цитаты валидируются перед отправкой
- Логи используют короткие стабильные поля: `err`, `operation`, `component`
- Секреты и сырые токены не логируются
