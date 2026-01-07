package reply

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) PinMessage(ctx context.Context, chatID int64, messageID int) error {
	msg := tgbotapi.PinChatMessageConfig{ChatID: chatID, MessageID: messageID, DisableNotification: true}

	_, err := t.bot.Send(msg)
	// Bug in SDK, Telegram API respond to Pin request just by bool value, not message
	if err != nil && !errors.Is(err, &json.UnmarshalTypeError{}) {
		return fmt.Errorf("failed to pin message: %w", err)
	}
	return nil
}
