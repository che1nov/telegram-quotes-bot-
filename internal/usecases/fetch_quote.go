package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"telegram-quotes-bot/internal/entities"
	"telegram-quotes-bot/internal/validators"
)

// QuoteAPI получает цитаты из внешнего источника.
type QuoteAPI interface {
	GetRandomQuote(ctx context.Context) (*entities.Quote, error)
}

// FetchQuoteService предоставляет методы для получения случайных цитат через API.
type FetchQuoteService struct {
	api    QuoteAPI // Интерфейс для взаимодействия с внешним API цитат
	logger *slog.Logger
}

// NewFetchQuoteService создаёт новый экземпляр FetchQuoteService.
// Принимает реализацию интерфейса QuoteAPI для получения цитат.
func NewFetchQuoteService(api QuoteAPI, logger *slog.Logger) *FetchQuoteService {
	if logger == nil {
		logger = slog.Default()
	}

	return &FetchQuoteService{
		api:    api,
		logger: logger.With("component", "fetch_quote_service"),
	}
}

// FetchQuote получает случайную цитату через API.
// Возвращает структуру Quote или ошибку, если не удалось получить цитату.
func (s *FetchQuoteService) FetchQuote(ctx context.Context) (*entities.Quote, error) {
	s.logger.InfoContext(ctx, "Получение цитаты", "operation", "fetch_quote")

	// Вызываем метод GetRandomQuote у переданного API для получения случайной цитаты
	quote, err := s.api.GetRandomQuote(ctx)
	if err != nil {
		// Если произошла ошибка при получении цитаты, возвращаем nil и сообщение об ошибке
		s.logger.ErrorContext(ctx, "Не удалось получить цитату", "operation", "fetch_quote", "err", err)
		return nil, fmt.Errorf("получить цитату: %w", err)
	}

	if quote == nil {
		err := fmt.Errorf("API вернул пустой результат")
		s.logger.ErrorContext(ctx, "Получена пустая цитата", "operation", "fetch_quote", "err", err)
		return nil, err
	}

	if err := validators.ValidateQuote(quote.Text, quote.Author); err != nil {
		s.logger.ErrorContext(ctx, "Получена невалидная цитата", "operation", "fetch_quote", "err", err)
		return nil, fmt.Errorf("получена невалидная цитата: %w", err)
	}

	s.logger.InfoContext(ctx, "Цитата получена", "operation", "fetch_quote")

	// Возвращаем полученную цитату
	return quote, nil
}
