package adapters

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestForismaticAPI_GetRandomQuote(t *testing.T) {
	// Простой тест создания API
	api := NewForismaticAPI(testLogger())
	if api == nil {
		t.Error("Expected API instance, got nil")
	}
	if api.client == nil {
		t.Error("Expected HTTP client, got nil")
	}

	// Тест с реальным API (может не работать в CI/CD)
	ctx := context.Background()
	quote, err := api.GetRandomQuote(ctx)

	// Если API недоступен, это нормально для тестов
	if err != nil {
		t.Logf("API недоступен (это нормально для тестов): %v", err)
		return
	}

	// Если API работает, проверяем результат
	if quote == nil {
		t.Error("Expected quote, got nil")
		return
	}

	if quote.Text == "" {
		t.Error("Expected non-empty quote text")
	}

	if quote.Author == "" {
		t.Error("Expected non-empty author")
	}

	t.Logf("Получена цитата: %q - %s", quote.Text, quote.Author)
}

func TestNewForismaticAPI(t *testing.T) {
	api := NewForismaticAPI(testLogger())
	if api == nil {
		t.Error("Expected API instance, got nil")
	}
	if api.client == nil {
		t.Error("Expected HTTP client, got nil")
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
