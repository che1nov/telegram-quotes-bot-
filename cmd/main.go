package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"telegram-quotes-bot/internal/adapters"
	"telegram-quotes-bot/internal/config"
	"telegram-quotes-bot/internal/usecases"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// setupLogger логгер
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	// Настройка логгера
	logger := setupLogger()

	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		logger.Warn("Файл .env не найден или не загружен")
	}

	// Создаем контекст с обработкой сигналов для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "err", err)
		os.Exit(1)
	}

	// Инициализация адаптеров
	quoteAPI := adapters.NewForismaticAPI(logger)
	telegramAdapter := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID, logger)

	// Инициализация сервисов
	fetchQuoteService := usecases.NewFetchQuoteService(quoteAPI, logger)
	sendQuoteService := usecases.NewSendQuoteService(telegramAdapter, logger)

	// Планировщик Cron
	c := cron.New()
	defer c.Stop()

	// Задача отправки цитат
	_, err = c.AddFunc("0 4,8,14,18 * * *", func() {
		taskCtx, taskCancel := context.WithTimeout(ctx, 45*time.Second)
		defer taskCancel()

		if err := sendQuote(taskCtx, fetchQuoteService, sendQuoteService); err != nil {
			logger.ErrorContext(taskCtx, "Ошибка отправки цитаты", "operation", "scheduled_quote", "err", err)
			return
		}

		logger.InfoContext(taskCtx, "Цитата успешно отправлена", "operation", "scheduled_quote")
	})
	if err != nil {
		logger.Error("Не удалось добавить cron-задачу", "err", err)
		os.Exit(1)
	}

	// Запуск планировщика
	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	// Отправка тестовой цитаты при запуске (если включено в конфигурации)
	if cfg.SendTestQuote {
		logger.Info("Отправка тестовой цитаты...")
		testCtx, testCancel := context.WithTimeout(ctx, 45*time.Second)
		defer testCancel()

		if err := sendQuote(testCtx, fetchQuoteService, sendQuoteService); err != nil {
			logger.ErrorContext(testCtx, "Ошибка отправки тестовой цитаты", "operation", "test_quote", "err", err)
		} else {
			logger.InfoContext(testCtx, "Тестовая цитата успешно отправлена", "operation", "test_quote")
		}
	} else {
		logger.Info("Отправка тестовой цитаты отключена в конфигурации")
	}

	// Ожидание сигнала завершения
	<-ctx.Done()
	logger.Info("Получен сигнал завершения. Останавливаем планировщик...")

	// Остановка планировщика
	stopCtx := c.Stop()
	<-stopCtx.Done()
	logger.Info("Планировщик остановлен. Программа завершена.")
}

func sendQuote(ctx context.Context, fetcher *usecases.FetchQuoteService, sender *usecases.SendQuoteService) error {
	quote, err := fetcher.FetchQuote(ctx)
	if err != nil {
		return err
	}

	if err := sender.SendQuote(ctx, quote); err != nil {
		return err
	}

	return nil
}
