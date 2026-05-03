package usecases

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"telegram-quotes-bot/internal/entities"
)

// mockQuoteAPI - мок для тестирования
type mockQuoteAPI struct {
	quote *entities.Quote
	err   error
}

func (m *mockQuoteAPI) GetRandomQuote(ctx context.Context) (*entities.Quote, error) {
	return m.quote, m.err
}

func TestFetchQuoteService_FetchQuote(t *testing.T) {
	tests := []struct {
		name          string
		mockAPI       *mockQuoteAPI
		expectedErr   bool
		expectedQuote *entities.Quote
	}{
		{
			name: "Success",
			mockAPI: &mockQuoteAPI{
				quote: &entities.Quote{
					Text:   "Test quote",
					Author: "Test Author",
				},
				err: nil,
			},
			expectedErr: false,
			expectedQuote: &entities.Quote{
				Text:   "Test quote",
				Author: "Test Author",
			},
		},
		{
			name: "API Error",
			mockAPI: &mockQuoteAPI{
				quote: nil,
				err:   errors.New("API error"),
			},
			expectedErr:   true,
			expectedQuote: nil,
		},
		{
			name: "Nil Quote",
			mockAPI: &mockQuoteAPI{
				quote: nil,
				err:   nil,
			},
			expectedErr:   true,
			expectedQuote: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewFetchQuoteService(tt.mockAPI, testLogger())
			ctx := context.Background()

			quote, err := service.FetchQuote(ctx)

			if tt.expectedErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if tt.expectedQuote != nil {
					if quote == nil {
						t.Error("Expected quote, got nil")
					} else {
						if quote.Text != tt.expectedQuote.Text {
							t.Errorf("Expected text %q, got %q", tt.expectedQuote.Text, quote.Text)
						}
						if quote.Author != tt.expectedQuote.Author {
							t.Errorf("Expected author %q, got %q", tt.expectedQuote.Author, quote.Author)
						}
					}
				}
			}
		})
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
