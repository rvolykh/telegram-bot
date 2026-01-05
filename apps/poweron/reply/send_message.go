package reply

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) SendMessage(ctx context.Context, chatID int64, message string) error {
	_, err := t.bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
