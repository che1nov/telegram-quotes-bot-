package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"telegram-quotes-bot/internal/validators"
)

// Config представляет конфигурацию приложения.
type Config struct {
	BotToken      string // Токен Telegram-бота
	ChatID        int64  // Идентификатор Telegram-канала
	SendTestQuote bool   // Отправлять ли тестовую цитату при запуске
}

// LoadConfig загружает конфигурацию из переменных окружения.
func LoadConfig(logger *slog.Logger) (*Config, error) {
	if logger == nil {
		logger = slog.Default()
	}

	// Чтение переменных окружения
	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")
	sendTestQuoteStr := os.Getenv("SEND_TEST_QUOTE")

	// Проверка наличия обязательных переменных
	if botToken == "" || chatIDStr == "" {
		logger.Error("Необходимые переменные окружения отсутствуют",
			"bot_token", maskToken(botToken), "chat_id", chatIDStr)
		return nil, errors.New("необходимые переменные окружения отсутствуют")
	}

	// Преобразование CHAT_ID в int64
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		logger.Error("Ошибка преобразования CHAT_ID в int64", "err", err)
		return nil, err
	}

	// Валидируем токен бота
	if err := validators.ValidateBotToken(botToken); err != nil {
		logger.Error("Невалидный токен бота", "err", err)
		return nil, err
	}

	// Валидируем ID чата
	if err := validators.ValidateChatID(chatID); err != nil {
		logger.Error("Невалидный ID чата", "err", err)
		return nil, err
	}

	// Парсим флаг отправки тестовой цитаты (по умолчанию true)
	sendTestQuote := true
	if sendTestQuoteStr != "" {
		if parsed, err := strconv.ParseBool(sendTestQuoteStr); err == nil {
			sendTestQuote = parsed
		} else {
			logger.Warn("Неверное значение SEND_TEST_QUOTE, используется значение по умолчанию (true)", "value", sendTestQuoteStr)
		}
	}

	// Возвращаем конфигурацию
	return &Config{
		BotToken:      botToken,
		ChatID:        chatID,
		SendTestQuote: sendTestQuote,
	}, nil
}

// maskToken маскирует токен для безопасного логирования
func maskToken(token string) string {
	if token == "" {
		return "***"
	}
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}
