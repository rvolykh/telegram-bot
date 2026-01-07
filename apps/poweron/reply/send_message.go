package reply

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) SendMessage(ctx context.Context, chatID int64, message string) (int, error) {
	result, err := t.bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		return 0, fmt.Errorf("telegram API has failed: %w", err)
	}
	return result.MessageID, nil
}
