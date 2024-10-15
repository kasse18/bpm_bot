package core

import (
	"bmp-tgbot/internal/db"
	"bmp-tgbot/internal/sdk"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
)

type TelegramBot struct {
	dbClient db.IClient
	bot      *tgbotapi.BotAPI
	logger   *zap.Logger
}

func NewTelegramBot(dbClient db.IClient, logger *zap.Logger) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(sdk.EnvToken))
	if err != nil {
		panic(err)
	}

	return &TelegramBot{
		dbClient: dbClient,
		bot:      bot,
		logger:   logger,
	}
}

func (r *TelegramBot) Start() {
	defer func() {
		r.logger.Info("stopped telegram bot")
	}()

	if err := tgbotapi.SetLogger(&TelegramLogger{r.logger}); err != nil {
		panic(err)
	}

	if os.Getenv(sdk.EnvDebug) != "" {
		r.bot.Debug = true
	}

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates := r.bot.GetUpdatesChan(config)

	var err error
	r.logger.Info("wait for updates")
	for update := range updates {
		err = r.process(&update)
		if err != nil {
			r.logger.Error("failed to process update", zap.Error(err))
		}
	}
}
