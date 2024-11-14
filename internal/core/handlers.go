package core

import (
	"bmp-tgbot/internal/sdk"
	"bmp-tgbot/internal/sdk/models"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (r *TelegramBot) handleStart(user *models.User) error {
	config := tgbotapi.NewMessage(user.ID, sdk.MessageStart)
	config.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(
				sdk.ButtonCasino,
			),
		),
	)

	if _, err := r.bot.Send(config); err != nil {
		return err
	}

	return nil
}

func (r *TelegramBot) handleCasino(ctx context.Context, user *models.User, update *tgbotapi.Update) error {
	if update.Message.ForwardFrom != nil {
		config := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: update.Message.MessageID,
			},
			Text:                  fmt.Sprintf("Саси жопу"),
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		}
		if _, err := r.bot.Send(config); err != nil {
			return err
		}
		return nil
	}
	if user.Balance <= 0 {
		config := tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: update.Message.MessageID}
		if _, err := r.bot.Send(config); err != nil {
			return err
		}
		return nil
	}
	val := update.Message.Dice.Value
	score := int64(0)
	switch val {
	case 64:
		score = 10
	case 1, 22, 43:
		score = 7
	case 16, 32, 48:
		score = 5
	default:
		score = -1
	}
	user.Balance = user.Balance + score
	if err := r.dbClient.UpdateUser(ctx, user); err != nil {
		r.logger.Error("failed to update user balance", zap.Error(err))
	}
	config := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           update.Message.Chat.ID,
			ReplyToMessageID: update.Message.MessageID,
		},
		Text:                  fmt.Sprintf("Your value: %v\nYour score: %v\nYour balance: %v", val, score, user.Balance),
		ParseMode:             "",
		Entities:              nil,
		DisableWebPagePreview: false,
	}
	//(BaseConfig{ChatID: update.Message.Chat.ID, ReplyToMessageID: update.Message.MessageID}, fmt.Sprintf("Your value: %v\nYour score: %v\nYour balance: %v", val, score, user.Balance))
	if _, err := r.bot.Send(config); err != nil {
		return err
	}
	return nil
}

func (r *TelegramBot) handleLeaderboard(ctx context.Context, update *tgbotapi.Update) error {
	board, err := r.dbClient.GetLeaderboard(ctx)
	if err != nil {
		r.logger.Error("failed to get leaderboard", zap.Error(err))
	}
	config := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           update.Message.Chat.ID,
			ReplyToMessageID: update.Message.MessageID,
		},
		Text:                  fmt.Sprintf("Leaderboard:\n%v", board),
		ParseMode:             "",
		Entities:              nil,
		DisableWebPagePreview: false,
	}
	//tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Leaderboard:\n%v", board))
	if _, err := r.bot.Send(config); err != nil {
		return err
	}
	return nil
}

func (r *TelegramBot) handleGetBalance(ctx context.Context, user *models.User, update *tgbotapi.Update) error {
	balance := user.Balance
	if balance > 0 {
		config := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: update.Message.MessageID,
			},
			Text:                  fmt.Sprintf("Баланс должен быть <= 0"),
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		}
		if _, err := r.bot.Send(config); err != nil {
			return err
		}
		return nil
	}
	user.Balance = user.Balance + 10
	if err := r.dbClient.UpdateUser(ctx, user); err != nil {
		r.logger.Error("failed to update user balance", zap.Error(err))
	}
	config := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           update.Message.Chat.ID,
			ReplyToMessageID: update.Message.MessageID,
		},
		Text:                  fmt.Sprintf("Баланс обновлен: %v", user.Balance),
		ParseMode:             "",
		Entities:              nil,
		DisableWebPagePreview: false,
	}

	if _, err := r.bot.Send(config); err != nil {
		return err
	}
	return nil
}

//
//func (r *TelegramBot) handleAll(ctx context.Context, update *tgbotapi.Update) {
//	update.Message.
//}
