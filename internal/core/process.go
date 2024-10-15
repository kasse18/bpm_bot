package core

import (
	"bmp-tgbot/internal/sdk"
	"bmp-tgbot/internal/sdk/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (r *TelegramBot) process(update *tgbotapi.Update) error {
	tgUser := r.getTgUserFromUpdate(update)
	fmt.Println(tgUser)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := r.dbClient.GetUser(ctx, tgUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = r.dbClient.CreateUser(ctx, tgUser)
		}

		if err != nil {
			return err
		}
	}

	switch {
	case update.Message.Command() != "":
		return r.handleCommand(ctx, tgUser, update)
	case update.Message.Dice != nil:
		return r.handleDice(ctx, tgUser, update)
	}

	return nil
}

func (r *TelegramBot) handleDice(ctx context.Context, user *models.User, update *tgbotapi.Update) error {
	return r.handleCasino(ctx, user, update)
}

func (r *TelegramBot) handleCommand(ctx context.Context, user *models.User, update *tgbotapi.Update) error {
	switch update.Message.Command() {
	case sdk.CommandStart:
		return r.handleStart(user)
	case sdk.CommandLeaderboard:
		return r.handleLeaderboard(ctx, update)
	}

	return nil
}
