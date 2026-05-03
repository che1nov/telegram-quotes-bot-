package usecases

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"unicode/utf8"

	"telegram-quotes-bot/internal/entities"
)

// mockTelegramSender - мок для тестирования
type mockTelegramSender struct {
	err     error
	message string
}

func (m *mockTelegramSender) SendMessage(ctx context.Context, message string) error {
	m.message = message
	return m.err
}

func TestSendQuoteService_SendQuote(t *testing.T) {
	tests := []struct {
		name        string
		mockSender  *mockTelegramSender
		quote       *entities.Quote
		expectedErr bool
	}{
		{
			name: "Success",
			mockSender: &mockTelegramSender{
				err: nil,
			},
			quote: &entities.Quote{
				Text:   "Test quote",
				Author: "Test Author",
			},
			expectedErr: false,
		},
		{
			name: "Send Error",
			mockSender: &mockTelegramSender{
				err: errors.New("send error"),
			},
			quote: &entities.Quote{
				Text:   "Test quote",
				Author: "Test Author",
			},
			expectedErr: true,
		},
		{
			name: "Nil Quote",
			mockSender: &mockTelegramSender{
				err: nil,
			},
			quote:       nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSendQuoteService(tt.mockSender, sendTestLogger())
			ctx := context.Background()

			err := service.SendQuote(ctx, tt.quote)

			if tt.expectedErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestSendQuoteService_FormatQuote_TruncatesUTF8Safely(t *testing.T) {
	service := NewSendQuoteService(&mockTelegramSender{}, sendTestLogger())
	quote := &entities.Quote{
		Text:   strings.Repeat("я", 201),
		Author: "Автор",
	}

	for i := 0; i < 5; i++ {
		message := service.FormatQuote(quote)
		if !utf8.ValidString(message) {
			t.Fatalf("expected valid UTF-8 message, got %q", message)
		}
	}
}

func sendTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
