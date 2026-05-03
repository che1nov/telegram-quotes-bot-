# Telegram Quotes Bot

[Русская версия](README.ru.md)

Telegram bot that automatically sends motivational quotes in Russian to a Telegram chat or channel on a schedule. The bot fetches random quotes from the Forismatic API and sends them through the Telegram Bot API.

## Features
- Receives random motivational quotes in Russian from the Forismatic API
- Sends quotes to Telegram on a cron schedule
- Can send a test quote on startup
- Handles graceful shutdown with OS signals
- Uses HTTP clients with timeouts
- Validates configuration and external quote data
- Uses structured logging with `log/slog`
- Keeps a simple layered structure with use cases and adapters

## Technologies
- **Language**: Go 1.23
- **APIs**: Telegram Bot API, Forismatic API
- **Scheduler**: robfig/cron
- **Architecture**: simple layered structure with dependency injection
- **Testing**: Go testing package
- **Containerization**: Docker
- **Task runner**: Taskfile

## Project Structure
```text
├── cmd/                   # Application entry point
├── internal/
│   ├── adapters/          # External service adapters
│   ├── config/            # Configuration loading
│   ├── entities/          # Domain entities
│   ├── usecases/          # Application scenarios
│   └── validators/        # Input validation
├── Dockerfile             # Container image
├── Taskfile.yml           # Task automation
└── README.ru.md           # Russian README
```

## Quick Start

### Prerequisites
- Go 1.23+
- Docker, optional
- Task, optional

### Setup
1. Clone the repository.
2. Create a `.env` file.
3. Set the Telegram bot token and target chat ID:
   ```bash
   BOT_TOKEN=123456789:your-token
   CHAT_ID=-1001234567890
   SEND_TEST_QUOTE=true
   ```

### Running
```bash
# Go
go run cmd/main.go

# Task
task run

# Docker
docker build -t telegram-quotes-bot .
docker run --env-file .env telegram-quotes-bot
```

### Development
```bash
# Run tests
task test

# Run tests with coverage
task test-coverage

# Build the app
task build

# Run go vet and gofmt
task lint

# Remove build artifacts
task clean
```

## Configuration
The bot reads configuration from environment variables:
- `BOT_TOKEN`: Telegram bot token from BotFather
- `CHAT_ID`: target Telegram chat or channel ID
- `SEND_TEST_QUOTE`: send a test quote on startup, `true` by default

## Schedule
The bot sends quotes at `04:00`, `08:00`, `14:00`, and `18:00` every day.

## API
The bot uses the Forismatic API to fetch Russian quotes:
- **Endpoint**: `http://api.forismatic.com/api/1.0/?method=getQuote&format=json&lang=ru`
- **Response fields**: `quoteText`, `quoteAuthor`
- **Language**: Russian

## Security
- Bot token is masked in logs
- HTTP clients have timeouts
- External quote data is validated before sending
- Logs use short stable fields such as `err`, `operation`, and `component`
- Secrets and raw tokens are not logged
