package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"telegram-quotes-bot/internal/entities"
)

// TelegramSender отправляет готовые сообщения в Telegram.
type TelegramSender interface {
	SendMessage(ctx context.Context, message string) error
}

// SendQuoteService предоставляет методы для отправки цитат в Telegram-канал.
type SendQuoteService struct {
	telegram TelegramSender // Интерфейс для отправки сообщений в Telegram
	logger   *slog.Logger
}

// NewSendQuoteService создаёт новый экземпляр SendQuoteService.
// Принимает интерфейс TelegramSender для отправки сообщений в Telegram.
func NewSendQuoteService(telegram TelegramSender, logger *slog.Logger) *SendQuoteService {
	if logger == nil {
		logger = slog.Default()
	}

	return &SendQuoteService{
		telegram: telegram,
		logger:   logger.With("component", "send_quote_service"),
	}
}

// SendQuote отправляет цитату в Telegram-канал.
// Форматирует цитату в удобочитаемый вид и отправляет её через TelegramSender.
// Возвращает ошибку, если отправка не удалась.
func (s *SendQuoteService) SendQuote(ctx context.Context, quote *entities.Quote) error {
	// Проверяем, что цитата не nil
	if quote == nil {
		return fmt.Errorf("цитата не может быть nil")
	}

	s.logger.InfoContext(ctx, "Отправка цитаты", "operation", "send_quote")

	// Форматируем цитату с красивым оформлением
	message := s.FormatQuote(quote)

	// Отправляем сформированное сообщение через TelegramSender
	err := s.telegram.SendMessage(ctx, message)
	if err != nil {
		// Если произошла ошибка при отправке, возвращаем её с описанием
		s.logger.ErrorContext(ctx, "Не удалось отправить цитату", "operation", "send_quote", "err", err)
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	s.logger.InfoContext(ctx, "Цитата отправлена", "operation", "send_quote")

	// Если всё прошло успешно, возвращаем nil
	return nil
}

// FormatQuote создает красиво отформатированное сообщение с цитатой (публичная функция для тестирования)
func (s *SendQuoteService) FormatQuote(quote *entities.Quote) string {
	// Ограничиваем длину цитаты для лучшего отображения
	text := truncateRunes(quote.Text, 200)

	// Выбираем случайный стиль форматирования
	styles := []func(string, string) string{
		s.formatStyle2,
		s.formatStyle3,
	}

	style := styles[rand.Intn(len(styles))]
	return style(text, quote.Author)
}

// formatStyle2 - Стиль с кавычками
func (s *SendQuoteService) formatStyle2(text, author string) string {
	emojis := []string{"💫", "✨", "🌟", "🎯", "🔥", "💡", "🌈", "🦋", "🌸", "🎪"}
	emoji := emojis[rand.Intn(len(emojis))]

	return fmt.Sprintf(
		"%s *Мудрая мысль*\n\n"+
			"❝ %s ❞\n\n"+
			"    — *%s* ✍️",
		emoji,
		text,
		author,
	)
}

func truncateRunes(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}

	return string(runes[:limit-3]) + "..."
}

// formatStyle3 - Стиль с разделителями
func (s *SendQuoteService) formatStyle3(text, author string) string {
	emojis := []string{"🌟", "💫", "✨", "🎯", "🔥", "💡", "🌈", "🦋", "🌸", "🎨"}
	emoji := emojis[rand.Intn(len(emojis))]

	return fmt.Sprintf(
		"%s *Вдохновение дня*\n\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"  %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"+
			"👤 *%s*",
		emoji,
		text,
		author,
	)
}
