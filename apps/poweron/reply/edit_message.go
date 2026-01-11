package reply

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) EditMessage(ctx context.Context, chatID int64, messageID int, text string) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to edit message %d: %w", messageID, err)
	}
	return nil
}
