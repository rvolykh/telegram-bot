package reply

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) SendMenu(ctx context.Context, chatID int64, message string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboard

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send menu: %w", err)
	}
	return nil
}
