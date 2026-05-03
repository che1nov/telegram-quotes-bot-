package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TelegramAdapter реализует интерфейс TelegramSender для отправки сообщений в Telegram.
type TelegramAdapter struct {
	token  string
	chatID int64
	client *http.Client
	logger *slog.Logger
}

// NewTelegramAdapter создаёт новый экземпляр TelegramAdapter.
// Принимает токен бота (botToken) и ID чата (chatID).
func NewTelegramAdapter(botToken string, chatID int64, logger *slog.Logger) *TelegramAdapter {
	if logger == nil {
		logger = slog.Default()
	}

	return &TelegramAdapter{
		token:  botToken,
		chatID: chatID,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: logger.With("component", "telegram_adapter"),
	}
}

// SendMessage отправляет текстовое сообщение в Telegram-чат.
// Принимает контекст (ctx) и текст сообщения (message).
// Возвращает ошибку, если сообщение не удалось отправить.
func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	form := url.Values{}
	form.Set("chat_id", fmt.Sprintf("%d", t.chatID))
	form.Set("text", message)

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("создать запрос Telegram: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	t.logger.InfoContext(ctx, "Отправка сообщения", "operation", "send_message", "chat_id", t.chatID)

	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.ErrorContext(ctx, "Ошибка запроса Telegram", "operation", "send_message", "chat_id", t.chatID, "err", err)
		return fmt.Errorf("отправить запрос Telegram: %w", err)
	}
	defer resp.Body.Close()

	var telegramResp struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		t.logger.ErrorContext(ctx, "Ошибка декодирования ответа Telegram", "operation", "send_message", "chat_id", t.chatID, "err", err)
		return fmt.Errorf("декодировать ответ Telegram: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !telegramResp.OK {
		err := fmt.Errorf("telegram status=%d description=%q", resp.StatusCode, telegramResp.Description)
		t.logger.ErrorContext(ctx, "Telegram вернул ошибку", "operation", "send_message", "chat_id", t.chatID, "err", err)
		return err
	}

	t.logger.InfoContext(ctx, "Сообщение отправлено", "operation", "send_message", "chat_id", t.chatID)
	return nil
}
