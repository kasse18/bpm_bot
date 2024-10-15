package main

import (
	"bmp-tgbot/internal/core"
	"bmp-tgbot/internal/db"
	"bmp-tgbot/internal/run"
	"bmp-tgbot/internal/sdk"
	"context"
	"os"
)

func main() {
	run.Init()

	postgresClient := db.NewPostgresClient(context.Background(), os.Getenv(sdk.EnvPostgres))

	logger := run.Logger.Named("telegram bot")

	bot := core.NewTelegramBot(postgresClient, logger)
	bot.Start()
}
